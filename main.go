package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	guestauth "github.com/TootieJin/pjsekai-overlay-APPEND/pkg/auth"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/pjsekaioverlay"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/scoremaker"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
	"github.com/fatih/color"
	"github.com/google/go-github/v57/github"
	"github.com/srinathh/gokilo/rawmode"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/sys/windows"
)

func checkUpdate() (string, string) {
	githubClient := github.NewClient(nil)
	release, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), "TootieJin", "pjsekai-overlay-APPEND")
	if err != nil {
		return "", ""
	}

	latestVersion := strings.TrimPrefix(release.GetTagName(), "v")
	if latestVersion == pjsekaioverlay.Version || pjsekaioverlay.Version == "0.0.0" {
		return "", ""
	}
	return latestVersion, release.GetHTMLURL()
}

func checkSubstrings(str []string, subs ...string) string {
	for _, s := range str {
		for _, sub := range subs {
			if strings.Contains(s, sub) {
				return sub
			}
		}
	}
	return ""
}

type rawDownloadOptions struct {
	ScorePath         string
	PublishedInfoPath string
	PublishedScoreID  string
	PublishedUserID   int
	SessionToken      string
	DataVersion       string
	AssetVersion      string
	AssetHash         string
	Region            string
	OutPath           string
	SaveChartFiles    bool
	GenerateRaw       bool
	ForceDerivative   bool
}

type rawDownloadResult struct {
	PublishedInfo     scoremaker.PublishedScoreInfo
	PublishedInfoRaws [][]byte
	USCRaw            []byte
	LevelData         sonolus.LevelData
	HasLevelData      bool
	PJSKData          []byte
	HasPJSKData       bool
	SavedPath         string
}

type cachedCredential struct {
	UserID       int
	Credential   string
	SessionToken string
	Source       string
}

func runRawChartDownload(opts rawDownloadOptions) (rawDownloadResult, error) {
	result := rawDownloadResult{}
	scorePath := opts.ScorePath
	region := opts.Region
	outPath := opts.OutPath
	rawArtifactBasePath := strings.TrimSuffix(outPath, filepath.Ext(outPath))
	if strings.TrimSpace(opts.PublishedScoreID) != "" {
		rawArtifactBasePath = filepath.Join(filepath.Dir(outPath), opts.PublishedScoreID+"_raw-data")
	}
	if strings.TrimSpace(outPath) == "" {
		return result, fmt.Errorf("生出力パスの指定が必要です (Raw output path is required)")
	}

	var listing = false
	derivativeRestricted := false
	var publishedInfo scoremaker.PublishedScoreInfo
	var publishedInfoRaws [][]byte
	var err error
	if strings.TrimSpace(opts.PublishedScoreID) != "" {
		// Ensure a valid guest account + session, registering one on first run
		appInfoEarly, earlyErr := scoremaker.FetchAppInfo(region, nil)
		if earlyErr != nil {
			return result, fmt.Errorf("アプリ情報の自動解決に失敗しました (Failed to auto-resolve app info) [%w]", earlyErr)
		}

		guestUserID, guestToken, guestAuthData, guestCF, guestErr := guestauth.EnsureGuest(region, appInfoEarly, "")
		if guestErr != nil {
			return result, fmt.Errorf("ゲストアカウントの確保に失敗しました (Failed to ensure guest account) [%w]", guestErr)
		}

		// Prefer explicit CLI flags over the guest account
		hadExplicitUser := opts.PublishedUserID > 0
		hadExplicitToken := strings.TrimSpace(opts.SessionToken) != ""
		if opts.PublishedUserID <= 0 {
			opts.PublishedUserID = guestUserID
		}
		if strings.TrimSpace(opts.SessionToken) == "" {
			opts.SessionToken = guestToken
		}
		// Fill auth_data from guest account if not explicitly set
		if strings.TrimSpace(opts.DataVersion) == "" {
			opts.DataVersion = guestAuthData.DataVersion
		}
		if strings.TrimSpace(opts.AssetVersion) == "" {
			opts.AssetVersion = guestAuthData.AssetVersion
		}
		if strings.TrimSpace(opts.AssetHash) == "" {
			opts.AssetHash = guestAuthData.AssetHash
		}
		_ = hadExplicitToken
		_ = guestCF
		appInfo, appErr := scoremaker.FetchAppInfo(region, nil)
		if appErr != nil {
			return result, fmt.Errorf("アプリ情報の自動解決に失敗しました (Failed to auto-resolve app info) [%w]", appErr)
		}
		auth := scoremaker.AuthData{
			DataVersion:  opts.DataVersion,
			AssetVersion: opts.AssetVersion,
			AssetHash:    opts.AssetHash,
		}

		// Build attempt list: explicit/guest first, plus credential-backed retry
		var guestCred string
		if entry := guestCF.FindForRegion(region); entry != nil && entry.UserID == opts.PublishedUserID {
			guestCred = entry.Credential
		}
		attempts := []cachedCredential{
			{
				UserID:       opts.PublishedUserID,
				SessionToken: opts.SessionToken,
				Credential:   guestCred,
				Source:       guestauth.DefaultCredentialFile,
			},
		}
		allowRetry := !hadExplicitUser || !hadExplicitToken

		var lastErr error
		for _, attempt := range attempts {
			publishedInfo, publishedInfoRaws, _, err = scoremaker.FetchPublishedScoreInfo(
				opts.PublishedScoreID,
				attempt.UserID,
				attempt.SessionToken,
				region,
				auth,
				appInfo,
			)
			if err != nil && strings.Contains(strings.ToLower(err.Error()), "session_error") {
				// refresh via guestauth
				refreshedUID, refreshedToken, refreshedAuth, refreshErr := guestauth.RefreshSession(region, appInfo, guestCF)
				if refreshErr == nil && strings.TrimSpace(refreshedToken) != "" {
					if refreshedAuth.DataVersion != "" {
						auth = refreshedAuth
					}
					publishedInfo, publishedInfoRaws, _, err = scoremaker.FetchPublishedScoreInfo(
						opts.PublishedScoreID,
						refreshedUID,
						refreshedToken,
						region,
						auth,
						appInfo,
					)
					if err == nil {
						attempt.SessionToken = refreshedToken
					}
				}
			}
			if err == nil {
				opts.PublishedUserID = attempt.UserID
				opts.SessionToken = attempt.SessionToken
				lastErr = nil
				break
			}
			lastErr = err
		}
		if lastErr != nil {
			if allowRetry && len(attempts) > 1 {
				return result, fmt.Errorf("%d回認証情報の入力を試みた後、公開された譜面情報の取得に失敗しました (Failed to fetch published chart info after %d credential attempts) [%w]", len(attempts), len(attempts), lastErr)
			}
			return result, fmt.Errorf("公開された譜面情報の取得に失敗しました (Failed to fetch published chart info [%w]", lastErr)
		}
		scorePath = publishedInfo.UserCustomMusicScorePath
	} else if strings.TrimSpace(opts.PublishedInfoPath) != "" {
		publishedInfo, err = scoremaker.LoadPublishedScoreInfoFile(opts.PublishedInfoPath)
		if err != nil {
			return result, fmt.Errorf("公開された譜面情報の読み込みに失敗しました (Failed to load published chart info) [%w]", err)
		}
		if strings.TrimSpace(scorePath) == "" {
			scorePath = publishedInfo.UserCustomMusicScorePath
		} else if scorePath != publishedInfo.UserCustomMusicScorePath {
			return result, fmt.Errorf(".pjsk譜面のパスが、公開された譜面の情報パスと一致しません (.pjsk chart path does not match published chart info path)")
		}
	}
	result.PublishedInfo = publishedInfo
	result.PublishedInfoRaws = publishedInfoRaws
	if !opts.ForceDerivative && strings.TrimSpace(result.PublishedInfo.UserCustomMusicScoreID) != "" && !result.PublishedInfo.IsDerivativeAllowed {
		derivativeRestricted = true
	}
	if strings.TrimSpace(scorePath) == "" {
		return result, fmt.Errorf(".pjsk譜面のパスは必須です (.pjsk chart path is required)")
	}
	appInfo, err := scoremaker.FetchAppInfo(region, nil)
	if err != nil {
		return result, fmt.Errorf("アプリ情報の自動解決に失敗しました (Failed to auto-resolve app info) [%w]", err)
	}

	banMaker, err := pjsekaioverlay.Listing("483473494141414141414141417733497751324149424145774934344A534C45456E7A6241434348784167476C6D6A354F733835674C737452444531694A6877644E6462714C356B68417A68793056624B556868545A6D4D6C2B79304735526D6C694E50307649346D316E7478763835474B72324957667A5A6339514256353863315541714634414141413D", strings.TrimSpace(fmt.Sprint(result.PublishedInfo.UserID)))
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return result, fmt.Errorf("failed to check listing: %w", err)
	} else if banMaker {
		listing = true
	}

	raw, err := scoremaker.DownloadRawChart(scorePath, scoremaker.DownloadConfig{
		Region:       region,
		AppVersion:   appInfo.AppVersion,
		AppHash:      appInfo.AppHash,
		IssueCookie:  "",
		SessionToken: opts.SessionToken,
	})
	if err != nil {
		return result, err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return result, err
	}

	uscRaw, err := scoremaker.ConvertPJSKToUSC(raw)
	if err != nil {
		return result, fmt.Errorf(".pjsk譜面を.uscに変換できませんでした (Failed to convert .pjsk chart to .usc) [%w]", err)
	}
	result.USCRaw = uscRaw
	if opts.SaveChartFiles && !derivativeRestricted {
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return result, err
		}
		if err := os.WriteFile(outPath, uscRaw, 0644); err != nil {
			return result, err
		}
	}
	result.SavedPath = outPath
	result.PJSKData = nil
	result.HasPJSKData = false

	if opts.GenerateRaw && (strings.TrimSpace(opts.PublishedInfoPath) != "" || strings.TrimSpace(opts.PublishedScoreID) != "") {
		sidecarPath := rawArtifactBasePath + ".published-score-info.json"
		sidecarRaw, err := json.MarshalIndent(publishedInfo, "", "  ")
		if err != nil {
			return result, err
		}
		if err := os.WriteFile(sidecarPath, sidecarRaw, 0644); err != nil {
			return result, err
		}
		for i, rawPayload := range publishedInfoRaws {
			if len(rawPayload) == 0 {
				continue
			}
			suffix := fmt.Sprintf(".published-score-api-response-%d.json", i)
			rawResponsePath := rawArtifactBasePath + suffix
			var responsePayload any
			if err := msgpack.Unmarshal(rawPayload, &responsePayload); err == nil {
				if normalized, err := json.Marshal(responsePayload); err == nil {
					rawPayload = normalized
				}
			}
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, rawPayload, "", "  "); err == nil {
				_ = os.WriteFile(rawResponsePath, pretty.Bytes(), 0644)
			} else {
				_ = os.WriteFile(rawResponsePath, rawPayload, 0644)
			}
		}
	}

	levelData, err := scoremaker.ConvertPJSKToNextSekaiLevelData(raw)
	if err != nil {
		fmt.Printf("- WARN: NextSekai変換に失敗。usc譜面は保存済みです (Conversion failed. .usc chart has been saved) [%s]\n", err.Error())
		return result, nil
	}
	if publishedInfo.MusicID > 0 {
		if susLevelData, susErr := scoremaker.ConvertSUSMusicIDToNextSekaiLevelData(region, publishedInfo.MusicID); susErr == nil {
			scoremaker.AppendEventArchetypesFromLevelData(&levelData, susLevelData)
		}
	}
	result.LevelData = levelData
	result.HasLevelData = true

	pjskData, err := scoremaker.ConvertRawPJSKToReadablePJSK(raw)
	if err != nil {
		return result, fmt.Errorf("生のpjsk譜面を、読み取り可能なpjskに変換できませんでした (Failed to convert raw .pjsk chart to readable .pjsk) [%w]", err)
	}
	result.PJSKData = pjskData
	result.HasPJSKData = true

	if opts.GenerateRaw && !listing && !derivativeRestricted {
		convertedPath := rawArtifactBasePath + ".next-sekai.json"
		convertedRaw, err := json.MarshalIndent(levelData, "", "  ")
		if err != nil {
			return result, err
		}
		if err := os.WriteFile(convertedPath, convertedRaw, 0644); err != nil {
			return result, err
		}

		convertedPath = rawArtifactBasePath + ".pjsk.json"
		if err := os.WriteFile(convertedPath, pjskData, 0644); err != nil {
			return result, err
		}
	}

	return result, nil
}

func locale() (string, error) {
	cmd := exec.Command("powershell", "-Command", "Get-WinSystemLocale | Select-Object -ExpandProperty Name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func langPackCheck() (string, error) {
	cmd := exec.Command("powershell", "-Command", "Get-InstalledLanguage")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func isAdminPerm(path string) bool {
	created := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return true
		}
		created = true
	}

	testFile := filepath.Join(path, ".test_access")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return true
	}

	// cleanup test file
	_ = os.Remove(testFile)

	if created {
		_ = os.Remove(path)
	}

	return false
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

func origMain(isOptionSpecified bool) {
	Title()

	var aviutlType int
	flag.IntVar(&aviutlType, "aviutl-type", 0, "AviUtlインスタンスを指定します。(Specify AviUtl instance.)\n'1': AviUtl\n'2': AviUtl ExEdit2")

	var skipAviutlModConfig bool
	flag.BoolVar(&skipAviutlModConfig, "skip-mod-config", false, "AviUtlの設定変更はスキップされます。(Skip modifying AviUtl configurations.)")

	var skipAviutlInstall bool
	flag.BoolVar(&skipAviutlInstall, "skip-obj-install", false, "AviUtlオブジェクトのインストールをスキップします。(Skip installing AviUtl objects.)")

	var skipAviutlScriptInstall bool
	flag.BoolVar(&skipAviutlScriptInstall, "skip-script-install", false, "AviUtlスクリプトのインストールをスキップします。(Skip installing AviUtl scripts.)")

	var noExplorerAutoOpen bool
	flag.BoolVar(&noExplorerAutoOpen, "no-explorer-auto-open", false, "出力先ディレクトリを自動で開くのを無効にします。(Disable auto-opening output directory in Explorer.)")

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./dist/_chartId_", "出力先ディレクトリを指定します。_chartId_ は譜面IDに置き換えられます。\nEnter the output path. _chartId_ will be replaced with the chart ID.")

	var chartInstance string
	flag.StringVar(&chartInstance, "instance", "", "サーバーインスタンス（またはソースURL）を指定します。\nSpecify the server instance (or source URL).")

	var customBG bool
	flag.BoolVar(&customBG, "custom-bg", false, "UntitledChartsでカスタム背景を使用する。(Use custom background in UntitledCharts.)")

	var scoreModeInt int
	flag.IntVar(&scoreModeInt, "score-mode", 1, "採点モードを指定します。(Specify scoring mode.)\n'1': デフォルト/Default\n'2': 大会モード/Tournament Mode (PERFECT = +3)")

	var teamPower float64
	flag.Float64Var(&teamPower, "power", 250000, "総合力を指定します。(Specify the team's power.)")

	var enUI bool
	flag.BoolVar(&enUI, "en-ui", false, "英語UIを使う(部分的な対応)\nUse English UI (Partial support)")

	var allFlick bool
	flag.BoolVar(&allFlick, "all-flick", false, "すべてのノーツをフリックとして扱います。(Treat all notes as flicks.)")

	var pjskRegion string
	flag.StringVar(&pjskRegion, "pjsk-region", "jp", "譜面メーカーの譜面を取得する地域を指定します。(Specify the region to fetch the Score Maker charts from.)")

	var forceDerivative bool
	flag.BoolVar(&forceDerivative, "force-derivative", false, "二次利用を禁止する譜面ファイルでも保存を強制します。(Force save Score Maker chart files that prohibits derivative use.)")

	// advanced flags

	var rawPublishedScoreID string
	flag.StringVar(&rawPublishedScoreID, "pjsk-published-score-id", "", "公開譜面IDを指定してuscを直接出力します。(Download a published score ID directly as usc.)")

	var rawOutPath string
	flag.StringVar(&rawOutPath, "pjsk-out", "./dist/pjsk-_pjskId_/chart-_pjskId_.usc", "公開譜面ダウンロードの出力先を指定します。_pjskId_ は公開譜面IDに置き換えられます。\nSpecify the output path for published score downloads. _pjskId_ will be replaced with the published score ID.")

	var generateRaw bool
	flag.BoolVar(&generateRaw, "pjsk-generate-raw", false, "生データJSONを出力します。(Write raw metadata/leveldata JSON files.)")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay-APPEND [オプション (Options)] [譜面ID (Chart ID)]")
		flag.PrintDefaults()
	}

	flag.Parse()

	var listing = false

	if strings.TrimSpace(rawPublishedScoreID) != "" || strings.TrimSpace(rawOutPath) != "./dist/pjsk-_pjskId_/chart-_pjskId_.usc" || generateRaw {
		fmt.Println("- pjsk譜面ダウンロード中 を取得してuscに変換中 (Downloading chart and converting to .usc)... ")
		rawOutPath = strings.ReplaceAll(rawOutPath, "_pjskId_", rawPublishedScoreID)
		result, err := runRawChartDownload(rawDownloadOptions{
			PublishedScoreID: rawPublishedScoreID,
			Region:           pjskRegion,
			OutPath:          rawOutPath,
			GenerateRaw:      generateRaw,
			SaveChartFiles:   false,
			ForceDerivative:  forceDerivative,
		})
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		banMaker, err := pjsekaioverlay.Listing("483473494141414141414141417733497751324149424145774934344A534C45456E7A6241434348784167476C6D6A354F733835674C737452444531694A6877644E6462714C356B68417A68793056624B556868545A6D4D6C2B79304735526D6C694E50307649346D316E7478763835474B72324957667A5A6339514256353863315541714634414141413D", strings.TrimSpace(fmt.Sprint(result.PublishedInfo.UserID)))
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		} else if banMaker {
			listing = true
		}

		derivativeRestricted := !forceDerivative && !result.PublishedInfo.IsDerivativeAllowed
		if derivativeRestricted {
			fmt.Println(color.HiYellowString("WARN: この譜面は二次利用を禁じています。譜面ファイルは保存できません。\nThis chart prohibits derivative use. Chart file cannot be saved."))
		} else if !listing {
			if err := os.MkdirAll(filepath.Dir(result.SavedPath), 0755); err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: 出力ディレクトリの作成に失敗しました。(Failed to create output directory.) [%s]", err.Error())))
				return
			}
			if err := os.WriteFile(result.SavedPath, result.USCRaw, 0644); err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: .usc出力の書き込みに失敗しました。(Failed to write .usc output.) [%s]", err.Error())))
				return
			}
			fmt.Printf("- 保存先 (Saved to): %s\n", color.CyanString(result.SavedPath))
			if result.HasPJSKData {
				readablePJSKPath := strings.TrimSuffix(result.SavedPath, filepath.Ext(result.SavedPath)) + ".pjsk.json"
				if err := os.WriteFile(readablePJSKPath, result.PJSKData, 0644); err != nil {
					fmt.Println(color.RedString(fmt.Sprintf("FAIL: failed to write readable pjsk output: %s", err.Error())))
					return
				}
				fmt.Printf("- 読めるPJSKデータ (Readable PJSK Data): %s\n", color.CyanString(readablePJSKPath))
			}
		}
		if generateRaw {
			rawArtifactBasePath := filepath.Join(filepath.Dir(rawOutPath), rawPublishedScoreID+"_raw-data")
			fmt.Printf("- 生データ情報 (Raw metadata): %s\n", color.CyanString(rawArtifactBasePath+".published-score-info.json"))
			if result.HasLevelData && !listing && !derivativeRestricted {
				fmt.Printf("- 生LevelData (Raw LevelData): %s\n", color.CyanString(rawArtifactBasePath+".next-sekai.json"))
			}
			if result.HasPJSKData && !listing && !derivativeRestricted {
				fmt.Printf("- 生PJSKデータ (Raw PJSK Data): %s\n", color.CyanString(rawArtifactBasePath+".pjsk.json"))
			}
		}
		if !noExplorerAutoOpen {
			cmd := exec.Command(`explorer`, `/select,`, filepath.Dir(rawOutPath))
			cmd.Run()
			time.Sleep(2000 * time.Millisecond)
		}
		return
	}

	latestVer, releaseURL := checkUpdate()
	if latestVer != "" {
		fmt.Printf(color.HiCyanString("新しいバージョンがリリースされています\nNew version released: v%s -> v%s\n"), pjsekaioverlay.Version, latestVer)
		fmt.Printf(color.HiCyanString("ダウンロード (Download Here) -> %s\n"), releaseURL)
		fmt.Println(color.RedString("\nFAIL: pjsekai-overlay-APPENDを最新バージョンに更新してください。\nUpdate pjsekai-overlay-APPEND to the latest version."))
		return
	}

	fmt.Printf("- 前提条件を確認中 (Checking prerequisites)... ")

	locale, err := locale()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	} else if locale != "ja-JP" {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL: お使いのシステムロケールが「日本語（日本）」に設定されていません。変更方法についてはWikiを参照してください。\nYour system locale is not set to \"Japanese (Japan)\". Refer to the wiki for how to change it.\n- System locale: %v", locale)))
		return
	}

	langPackCheck, err := langPackCheck()
	if err != nil {
		fmt.Println(color.HiYellowString(fmt.Sprintf("WARN: 言語パックを確認できません。(Unable to check language pack.)\n%s", err.Error())))
		// return
		// (temporary pass)
	} else if !strings.Contains(langPackCheck, "ja-JP") {
		fmt.Println(color.RedString("\nFAIL: 日本語言語パックがインストールされていません。変更方法についてはWikiを参照してください。\nJapanese language pack is not installed. Refer to the wiki for how to install it."))
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	if isAdminPerm(cwd) {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL: ディレクトリには管理者権限が必要です。pjsekai-overlay-APPENDを「C:\\」または別の場所に移動してください。\nYour directory requires administrative permissions. Please move pjsekai-overlay-APPEND to \"C:\\\" or somewhere else.\n\n出力先ディレクトリ (Output path): %s", cwd)))
		return
	}
	if !isASCII(cwd) {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL: ディレクトリに非ASCII文字が含まれています。pjsekai-overlay-APPENDを「C:\\」または別の場所に移動してください。\nYour directory contains non-ASCII characters. Please move pjsekai-overlay-APPEND to \"C:\\\" or somewhere else.\n\n出力先ディレクトリ (Output path): %s", cwd)))
		return
	}

	mappingName, mapping := pjsekaioverlay.SetOverlayDefault()

	if len(mapping) != 25 {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL:「default.ini」ファイルのデータに異常があります。「default.ini」ファイルを削除し、プログラムを再起動して再生成してください。\nAbnormal \"default.ini\" data. Please regenerate by deleting the \"default.ini\" file and reopen the program.\n- Mapping count: %v != 25", len(mapping))))
		return
	}

	var mappingFloat64 []float64
	for _, v := range mapping {
		v = strings.TrimRightFunc(v, func(r rune) bool {
			return strings.HasSuffix(string(r), "+") || strings.HasSuffix(string(r), "-") || strings.HasSuffix(string(r), ".") || (r < '0' || r > '9')
		})
		mappingFloat64 = append(mappingFloat64, func() float64 {
			f, _ := strconv.ParseFloat(v, 64)
			return f
		}())
	}

	var float64Pointer = func(f float64) *float64 {
		return &f
	}

	var inRange = map[string]bool{
		// Root
		"offset":      mappingFloat64[0] >= -99999.99 && mappingFloat64[0] <= 99999.99,
		"cache":       mappingFloat64[1] == 0 || mappingFloat64[1] == 1,
		"text_lang":   mappingFloat64[2] == 0 || mappingFloat64[2] == 1,
		"watermark":   mappingFloat64[3] == 0 || mappingFloat64[3] == 1,
		"detail_stat": mappingFloat64[4] == 0 || mappingFloat64[4] == 1,
		// Life
		"life":        mappingFloat64[5] >= 0 && mappingFloat64[5] <= 9999 && math.Mod(mappingFloat64[5], 1.0) == 0,
		"life_skill":  mappingFloat64[6] >= 0 && mappingFloat64[6] <= 2 && math.Mod(mappingFloat64[6], 1.0) == 0,
		"overflow":    mappingFloat64[7] == 0 || mappingFloat64[7] == 1,
		"lead_zero":   mappingFloat64[8] == 0 || mappingFloat64[8] == 1,
		"unlock_life": mappingFloat64[9] == 0 || mappingFloat64[9] == 1,
		// Score
		"min_digit":   mappingFloat64[10] >= 1 && mappingFloat64[10] <= 99 && math.Mod(mappingFloat64[10], 1.0) == 0,
		"score_skill": mappingFloat64[11] >= 0 && mappingFloat64[11] <= 2 && math.Mod(mappingFloat64[11], 1.0) == 0,
		"score_speed": mappingFloat64[12] >= 0,
		"anim_score":  mappingFloat64[13] == 0 || mappingFloat64[13] == 1,
		"wds_anim":    mappingFloat64[14] == 0 || mappingFloat64[14] == 1,
		// Combo
		"ap":               mappingFloat64[15] == 0 || mappingFloat64[15] == 1,
		"tag":              mappingFloat64[16] == 0 || mappingFloat64[16] == 1,
		"last_digit":       mappingFloat64[17] >= 0 && math.Mod(mappingFloat64[17], 1.0) == 0,
		"combo_speed":      mappingFloat64[18] >= 0,
		"combo_burst":      mappingFloat64[19] == 0 || mappingFloat64[19] == 1,
		"achievement_rate": float64Pointer(mappingFloat64[20]) != nil,
		// Skill
		"skill_speed": mappingFloat64[21] >= 0,
		"skill_cache": mappingFloat64[22] == 0 || mappingFloat64[22] == 1,
		// Judgement
		"judge":       mappingFloat64[23] >= 1 && mappingFloat64[23] <= 10 && math.Mod(mappingFloat64[23], 1.0) == 0,
		"judge_speed": mappingFloat64[24] >= 0,
	}

	var mappingErr []string
	for i := range mapping {
		inRangeBool := inRange[mappingName[i]]
		if !inRangeBool {
			mappingErr = append(mappingErr, mappingName[i], fmt.Sprintf("%v", mapping[i]))
		}
	}

	if mappingErr != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:値の1つが範囲外です。Wikiで指定されている範囲内に収まるよう、値を調整してください。\nOne of the value is out of range. Please adjust the value so that it's within range specified in the Wiki.\n- Mapping out of range: %s", mappingErr)))
		return
	}

	var mappingStr []string
	for _, v := range mapping {
		mappingStr = append(mappingStr, fmt.Sprintf("%v", v))
	}

	fmt.Println(color.GreenString("OK"))

	var aviutlPath, aviutlProcess, aviutlName string

	switch aviutlType {
	case 1:
		aviutlProcess = "aviutl.exe"
		aviutlName = "AviUtl"
		aviutlPath, _, _ = pjsekaioverlay.DetectAviUtl()
	case 2:
		aviutlProcess = "aviutl2.exe"
		aviutlName = "AviUtl ExEdit2"
		aviutlPath, _ = filepath.Abs("C:\\ProgramData\\aviutl2")
	default:
		aviutlPath, aviutlProcess, aviutlName = pjsekaioverlay.DetectAviUtl()
		if aviutlProcess != "" {
			fmt.Printf("Instance (auto-detected): %s\n", color.GreenString(aviutlName))
		}

		if aviutlProcess == "" {
			fmt.Print("ファイルを生成するAviUtlインスタンスを選択してください。\nChoose AviUtl instance to generate files.\n\n'1': AviUtl\n'2': AviUtl ExEdit2\n> ")
			before, _ := rawmode.Enable()
			tmpAviutlByte, _ := bufio.NewReader(os.Stdin).ReadByte()
			tmpAviutl := string(tmpAviutlByte)
			rawmode.Restore(before)
			switch tmpAviutl {
			default:
				aviutlProcess = ""
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpAviutl))
				fmt.Println(color.RedString("FAIL: AviUtlインスタンスが選択されていません。\nAviUtl instance not selected."))
				return
			case "1":
				aviutlProcess = "aviutl.exe"
				aviutlName = "AviUtl"
				aviutlPath, _, _ = pjsekaioverlay.DetectAviUtl()
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpAviutl))
				fmt.Println(color.GreenString("Instance: AviUtl"))
			case "2":
				aviutlProcess = "aviutl2.exe"
				aviutlName = "AviUtl ExEdit2"
				aviutlPath, _ = filepath.Abs("C:\\ProgramData\\aviutl2")
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpAviutl))
				fmt.Println(color.GreenString("Instance: AviUtl ExEdit2"))
			}
		}
	}

	var successInstall = false
	if !skipAviutlModConfig {
		success := pjsekaioverlay.ModifyAviUtlConfig(aviutlPath, aviutlProcess)
		if success {
			fmt.Println(color.GreenString(aviutlName + "の設定変更が正常に完了しました。(" + aviutlName + " configurations successfully modified.)"))
			successInstall = true
		}
	}
	if !skipAviutlInstall {
		success := pjsekaioverlay.TryInstallObject(aviutlPath, aviutlProcess, mappingStr)
		if success {
			fmt.Println(color.GreenString(aviutlName + "オブジェクトのインストールに成功しました。(" + aviutlName + " object successfully installed.)"))
			successInstall = true
		}
	}
	if !skipAviutlScriptInstall {
		success := pjsekaioverlay.TryInstallScript(aviutlPath, aviutlProcess)
		if success {
			fmt.Println(color.GreenString(aviutlName + "依存関係スクリプトのインストールに成功しました。(" + aviutlName + " dependency scripts successfully installed.)"))
			successInstall = true
		}
	}
	if successInstall {
		fmt.Println(color.HiYellowString("変更を適用するには、" + aviutlName + "を再起動してください。(Please restart " + aviutlName + " to apply changes.)\n"))
	}

	Tips()

	var chartId string
	if flag.Arg(0) != "" {
		chartId = flag.Arg(0)
		fmt.Printf("譜面ID (Chart ID): %s\n", color.GreenString(chartId))
	} else {
		fmt.Print("譜面IDを接頭辞込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'pjsk-': 譜面メーカー / Score Maker\n'sekai-best-': Sekai Viewer (sonolus.sekai.best)\n'sss-': Sbuga's Sonolus Server (sonolus.sbuga.com)\n'chcy-': Chart Cyanvas\n'ptlv-': Potato Leaves (ptlv.milkbun.org)\n'utsk-': Untitled Sekai (us.pim4n-net.com)\n'UnCh-': UntitledCharts (untitledcharts.com)\n'lalo-': laoloser's server (sonolus.laoloser.com)\n'skyra-': osciris's server (Skyra)\n'sync-': Local Server (ScoreSync + ScoreSync Modern)\n'custom-': Custom Server (Source URL)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	if after, ok := strings.CutPrefix(chartId, "pjsk-"); ok {
		publishedScoreID := strings.TrimSpace(after)
		if publishedScoreID == "" {
			fmt.Println(color.RedString("FAIL: pjsk-の後に公開譜面IDを指定してください。(Please provide a published chart ID after pjsk-.)"))
			return
		}

		fmt.Printf("- 譜面を取得中 (Getting chart): %s\n", color.CyanString(publishedScoreID))

		formattedOutDir := filepath.Join(cwd, strings.ReplaceAll(outDir, "_chartId_", chartId))
		resultDir := filepath.Dir(formattedOutDir) + "\\" + chartId
		uscPath := filepath.Join(formattedOutDir, "chart-"+publishedScoreID+".usc")

		fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(resultDir))
		fmt.Print("- 譜面を取得してuscに変換中 (Downloading chart and converting to .usc)... ")
		rawResult, err := runRawChartDownload(rawDownloadOptions{
			PublishedScoreID: publishedScoreID,
			Region:           pjskRegion,
			OutPath:          uscPath,
			GenerateRaw:      generateRaw,
			SaveChartFiles:   false,
			ForceDerivative:  forceDerivative,
		})
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		publishedInfo := rawResult.PublishedInfo
		if !rawResult.HasLevelData {
			fmt.Println(color.RedString("FAIL: 譜面ファイルを変換できませんでした。(Unable to convert the chart file.)"))
			return
		}

		banMaker, err := pjsekaioverlay.Listing("483473494141414141414141417733497751324149424145774934344A534C45456E7A6241434348784167476C6D6A354F733835674C737452444531694A6877644E6462714C356B68417A68793056624B556868545A6D4D6C2B79304735526D6C694E50307649346D316E7478763835474B72324957667A5A6339514256353863315541714634414141413D", strings.TrimSpace(fmt.Sprint(publishedInfo.UserID)))
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		} else if banMaker {
			listing = true
		}

		derivativeRestricted := !forceDerivative && !publishedInfo.IsDerivativeAllowed
		if derivativeRestricted {
			fmt.Println(color.HiYellowString("WARN: この譜面は二次利用を禁じています。譜面ファイルは保存できません。\nThis chart prohibits derivative use. Chart file cannot be saved."))
		} else if !listing && rawResult.HasPJSKData {
			if err := os.MkdirAll(filepath.Dir(rawResult.SavedPath), 0755); err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: 出力ディレクトリの作成に失敗しました。(Failed to create output directory.) [%s]", err.Error())))
				return
			}
			if err := os.WriteFile(rawResult.SavedPath, rawResult.USCRaw, 0644); err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: .usc出力の書き込みに失敗しました。(Failed to write .usc output.) [%s]", err.Error())))
				return
			}
			fmt.Printf("- 保存先 (Saved to): %s\n", color.CyanString(rawResult.SavedPath))
			readablePJSKPath := strings.TrimSuffix(rawResult.SavedPath, filepath.Ext(rawResult.SavedPath)) + ".pjsk.json"
			if err := os.WriteFile(readablePJSKPath, rawResult.PJSKData, 0644); err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: 読み取り可能な.pjsk出力の書き込みに失敗しました。(Failed to write readable .pjsk output.) [%s]", err.Error())))
				return
			}
			fmt.Printf("- 読めるPJSKデータ (Readable PJSK Data): %s\n", color.CyanString(readablePJSKPath))
		}
		levelData := rawResult.LevelData

		meta, metaErr := scoremaker.FetchMusicMetadata(pjskRegion, publishedInfo.MusicID, nil)
		if metaErr != nil {
			fmt.Println(color.HiYellowString(fmt.Sprintf("WARN: 楽曲メタデータを取得できませんでした (Failed to fetch music metadata.) [%s]", metaErr.Error())))
		}

		title := strings.TrimSpace(meta.Title)
		if title == "" {
			title = strings.TrimSpace(publishedInfo.Title)
		}
		if title == "" {
			title = publishedScoreID
		}
		publishedScoreTitle := strings.TrimSpace(publishedInfo.Title)
		if publishedScoreTitle == "" {
			publishedScoreTitle = title
		}
		publishedUserName := strings.TrimSpace(publishedInfo.UserName)
		if publishedUserName == "" {
			publishedUserName = "-"
		}

		titlev1 := strings.TrimSpace(publishedInfo.Title)
		if titlev1 == "" {
			titlev1 = publishedScoreID
		}

		lyricist := strings.TrimSpace(meta.Lyricist)
		composer := strings.TrimSpace(meta.Composer)
		arranger := strings.TrimSpace(meta.Arranger)
		if lyricist == "" {
			lyricist = "-"
		}
		if composer == "" {
			composer = "-"
		}
		if arranger == "" {
			arranger = "-"
		}

		extra := strings.TrimSpace(meta.Label)
		if extra == "" {
			if enUI {
				extra = "【Additional Info】"
			} else {
				extra = "【追加情報】"
			}
		}

		fmt.Printf("%s (%s) - %s (Lv. %s)\n",
			color.CyanString(title),
			color.CyanString(publishedInfo.Title),
			color.CyanString(publishedUserName),
			color.MagentaString(strconv.Itoa(publishedInfo.PlayLevel)),
		)

		fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
		jacketURL := fmt.Sprintf("https://storage.sekai.best/sekai-%s-assets/music/jacket/jacket_s_%03d/jacket_s_%03d.png", pjskRegion, publishedInfo.MusicID, publishedInfo.MusicID)
		jacketSource := pjsekaioverlay.Source{Id: "pjsk", Name: "PJSK Published", Host: "storage.sekai.best"}
		jacketLevel := sonolus.LevelInfo{Cover: sonolus.SRL{Url: jacketURL}}
		err = pjsekaioverlay.DownloadJacket(jacketSource, jacketLevel, formattedOutDir)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		fmt.Print("- ローカルで背景を生成中 - お待ちください (Generating background locally - please wait)... ")
		err = pjsekaioverlay.DownloadBackground(jacketSource, jacketLevel, formattedOutDir, chartId, "-v 1", false, true)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		err = pjsekaioverlay.DownloadBackground(jacketSource, jacketLevel, formattedOutDir, chartId, "-v 3", false, true)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		var scoreMode string
		switch scoreModeInt {
		default:
			scoreMode = "default"
		case 2:
			scoreMode = "tournament"
		}
		if !isOptionSpecified {
			fmt.Print("\n採点モードを選択してください。(Choose scoring mode.)\n'1': デフォルト/Default\n'2': 大会モード/Tournament Mode (PERFECT = +3)\n> ")
			before, _ := rawmode.Enable()
			tmpScoreModeByte, _ := bufio.NewReader(os.Stdin).ReadByte()
			tmpScoreMode := string(tmpScoreModeByte)
			rawmode.Restore(before)
			switch tmpScoreMode {
			default:
				scoreMode = "default"
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpScoreMode))
				fmt.Println(color.GreenString("Score Mode: デフォルト/Default"))
			case "2":
				scoreMode = "tournament"
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpScoreMode))
				fmt.Println(color.GreenString("Score Mode: 大会/Tournament"))
			}
		}

		if !isOptionSpecified && scoreMode == "default" {
			fmt.Print("\n総合力を指定してください。 (Input your team power.)\n\n- 小数と科学的記数法が使える (Accepts decimals & scientific notation)\n- おすすめ (Recommended): 250000 - 300000\n- 例 (Example): 1234567; 1e+20; -300000\n> ")
			var tmpTeamPower string
			fmt.Scanln(&tmpTeamPower)
			if tmpTeamPower == "" {
				tmpTeamPower = "250000"
			}
			teamPower, err = strconv.ParseFloat(tmpTeamPower, 64)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		}

		fmt.Print("- スコアを計算中 (Calculating score)... ")
		rating := publishedInfo.PlayLevel
		calcLevel := sonolus.LevelInfo{Rating: rating}
		scoreData := pjsekaioverlay.CalculateScore(calcLevel, levelData, teamPower, scoreMode, allFlick, listing)
		fmt.Println(color.GreenString("OK"))

		if !isOptionSpecified {
			fmt.Print("\n英語UIを使う？（部分的な対応）[y/n]\nUse English UI? (Partial support) [y/n]\n> ")
			before, _ := rawmode.Enable()
			tmpEnableENByte, _ := bufio.NewReader(os.Stdin).ReadByte()
			tmpEnableEN := string(tmpEnableENByte)
			rawmode.Restore(before)
			if tmpEnableEN == "Y" || tmpEnableEN == "y" {
				enUI = true
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableEN))
				fmt.Println(color.GreenString("TOGGLE: ON"))
			} else {
				enUI = false
				fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableEN))
				fmt.Println(color.RedString("TOGGLE: OFF"))
			}
		}

		assets := filepath.Join(cwd, "assets")

		fmt.Print("- pedファイルを生成中 (Generating ped file)... ")
		err = pjsekaioverlay.WritePedFile(scoreData, assets, filepath.Join(formattedOutDir, "data.ped"), calcLevel, levelData, scoreMode, enUI, listing)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		var exoType = "exo"
		if aviutlProcess == "aviutl2.exe" {
			exoType = "alias(.object)"
		}

		fmt.Printf("- %sファイルを生成中 (Generating %s file)... ", exoType, exoType)
		difficulty := strings.ToUpper(strings.TrimSpace(publishedInfo.MusicDifficultyType))
		if difficulty == "" {
			difficulty = "APPEND"
		}

		description := []string{
			fmt.Sprintf("作詞：%s    作曲：%s    編曲：%s", lyricist, composer, arranger),
			"Vo：-",
			fmt.Sprintf("%s", publishedScoreTitle),
			fmt.Sprintf("%s", publishedUserName),
		}
		descriptionv1 := []string{
			fmt.Sprintf("作詞：%s    作曲：%s    編曲：%s", lyricist, composer, arranger),
			fmt.Sprintf("歌：-    譜面制作：%s", publishedUserName),
		}
		exFile := "tournament-mode.png"
		exFileOpacity := "100.0"
		if enUI {
			description = []string{
				fmt.Sprintf("Lyrics: %s    Music: %s    Arrangement: %s", lyricist, composer, arranger),
				"Vocals: -",
				fmt.Sprintf("%s", publishedScoreTitle),
				fmt.Sprintf("%s", publishedUserName),
			}
			descriptionv1 = []string{
				fmt.Sprintf("Lyrics: %s    Music: %s    Arrangement: %s", lyricist, composer, arranger),
				fmt.Sprintf("Vocals: -    Chart Design: %s", publishedUserName),
			}
			exFile = "tournament-mode-en.png"
		}
		if scoreMode == "tournament" {
			exFileOpacity = "0.0"
		}

		if aviutlProcess == "aviutl.exe" {
			err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, title, titlev1, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
		} else {
			err = pjsekaioverlay.WriteAliasFiles(assets, formattedOutDir, title, titlev1, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
		}
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		fmt.Println(color.GreenString("OK"))

		message := fmt.Sprintf("\n全ての処理が完了しました！READMEの規約を確認した上で、%sファイルを%sにインポートして下さい。\nExecution complete! Please import the %s file into %s after reviewing the README Terms of Use.", exoType, aviutlName, exoType, aviutlName)
		fmt.Println(color.GreenString(message))

		if !isOptionSpecified || !noExplorerAutoOpen {
			cmd := exec.Command(`explorer`, `/select,`, resultDir)
			cmd.Run()
			time.Sleep(2000 * time.Millisecond)
		}
		return
	}

	// Instance section
	if chartInstance == "" && strings.HasPrefix(chartId, "chcy-") {
		fmt.Printf("\nChart Cyanvasインスタンスを選択してください。(Please choose Chart Cyanvas instance.)\n\n[インスタンス一覧 - List of instance(s)]\n'0': アーカイブ/Archive - cc.milkbun.org\n'1': 分岐サーバー/Offshoot server - chart-cyanvas.com\n> ")
		var chartInput string
		fmt.Scanln(&chartInput)

		chartInstance = chartInput
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartInput))

	} else if chartInstance == "" && strings.HasPrefix(chartId, "custom-") {
		fmt.Printf("\nサーバーのソースURLを入力してください。(Enter the server source URL.)\n> ")
		var chartInput string
		fmt.Scanln(&chartInput)
		chartInput = strings.TrimPrefix(chartInput, "http://")
		chartInput = strings.TrimPrefix(chartInput, "https://")

		chartInstance = chartInput
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartInput))
	}

	var chartSource pjsekaioverlay.Source
	if strings.HasPrefix(chartId, "sync") {
		chartSource, err = pjsekaioverlay.DetectLocalChartSource()
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
		if strings.Contains(chartId, "-") {
			parts := strings.SplitN(chartId, "-", 2)
			if len(parts) == 2 {
				chartId = parts[1]
			}
		} else {
			fmt.Print("ローカルサーバーの譜面を入力してください。(Enter chart ID for the local server.)\n> ")
			fmt.Scanln(&chartId)
		}
	} else {
		chartSource, err = pjsekaioverlay.DetectChartSource(chartId, chartInstance)
		if err != nil {
			fmt.Println(color.RedString("FAIL: 譜面が見つかりません。接頭辞も込め、正しい譜面IDを入力して下さい。\nChart not found. Please enter the correct chart ID including the prefix."))
			return
		}

		switch chartSource.Status {
		case 1:
			fmt.Printf(color.RedString("FAIL: %sは対応されなくなりました。ご利用ありがとうございました。\n%s is no longer supported. Thank you for using the service.\n"), chartSource.Name, chartSource.Name)
			return
		case 2:
			fmt.Printf(color.HiYellowString("WARN: %sは現在開発中であり、正常に動作しない可能性があります。\n%s is currently in development and may not work.\n"), chartSource.Name, chartSource.Name)
		}
	}

	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())

	var chart sonolus.LevelInfo
	prefixTrim := checkSubstrings([]string{chartId}, "lalo-", "skyra-", "custom-")
	chart, err = pjsekaioverlay.FetchChart(chartSource, strings.TrimPrefix(chartId, prefixTrim))

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	// Additional BG
	chartCCv1, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?c_background=v1")
	// chartUNv3, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=v3")
	// chartUNv1, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=v1")
	chartUNv1def, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=default_or_v1")

	if chart.Engine.Version != 13 {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL (ver.%d): エンジンのバージョンが古い。pjsekai-overlay-APPENDを最新版に更新してください。\nUnsupported engine version. Please update pjsekai-overlay-APPEND to latest version.", chart.Engine.Version)))
		return
	}

	banSource, err := pjsekaioverlay.Listing("4834734941414141414141414177584255524B4145425146304231357A7969616C74423347304443544E46775453322F63784C77394A556F356734524D394A776F34666D6130456F454C3765744E654B484C5A63614A4B4850767A436939576E4D737A73316179636B5534474F55397371646D586E43326A58514966667052774B57396341414141", chart.Source)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	} else if banSource {
		listing = true
	}
	banList, err := pjsekaioverlay.Listing("4834734941414141414141414178584D32773241494177463049327344314C4645667832675A59416B696759755562484E3534427A6761636453614B71614B4A43647574642F57584B786B2B6F33486C6F4C55554A4C2B6B54494D3474734C61427862726A50672B54454537566A5857634476534A512B7035503350384F49446F7A33526A3130414141413D", chart.Author)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	} else if banList {
		listing = true
	}

	fmt.Println(color.GreenString("OK"))
	fmt.Printf("  %s / %s - %s (Lv. %s)\n",
		color.CyanString(chart.Title),
		color.CyanString(chart.Artists),
		color.CyanString(chart.Author),
		color.MagentaString(strconv.Itoa(chart.Rating)),
	)

	fmt.Printf("- exeのパスを取得中 (Getting executable path)... ")
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.ReplaceAll(outDir, "_chartId_", chartId))
	resultDir := filepath.Dir(formattedOutDir) + "\\" + chartId

	fmt.Println(color.GreenString("OK"))
	fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(resultDir))

	fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
	err = pjsekaioverlay.DownloadJacket(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	/*
		fmt.Print("- 音声のプレビューをダウンロード中 (Downloading preview audio)... ")
		err = pjsekaioverlay.DownloadPreview(chartSource, chart, formattedOutDir)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		fmt.Println(color.GreenString("OK"))
	*/

	if !isOptionSpecified && (chartSource.Id == "untitledcharts" || chartSource.Id == "skyra") {
		fmt.Print("\nカスタム背景を使用しますか？（デフォルトを使用するには「n」）[y/n]\nUse custom background? ('n' to use default) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpCustomBGByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpCustomBG := string(tmpCustomBGByte)
		rawmode.Restore(before)
		if tmpCustomBG == "Y" || tmpCustomBG == "y" {
			customBG = true
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpCustomBG))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			customBG = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpCustomBG))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}

	const localGenerate = true
	if customBG {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "", customBG, !localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		if chartSource.Id == "untitledcharts" {
			err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1def, formattedOutDir, chartId+"?levelbg=default_or_v1", "", customBG, !localGenerate)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		} else {
			err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId+"/", "", customBG, !localGenerate)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		}
	} else if chartSource.Id == "untitledcharts" {
		/*
			fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

			err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv3, formattedOutDir, chartId+"?levelbg=v3", "", customBG, !localGenerate)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}

			err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1, formattedOutDir, chartId+"?levelbg=v1", "", customBG, !localGenerate)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}

			// their background is saturated
		*/
		fmt.Print("- ローカルで背景を生成中 - お待ちください (Generating background locally - please wait)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 1", customBG, localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 3", customBG, localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else if chartSource.Id == "chart_cyanvas" && chartSource.Name != "Chart Cyanvas Archive" {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "", customBG, !localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartCCv1, formattedOutDir, chartId+"?c_background=v1", "", customBG, !localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else {
		fmt.Print("- ローカルで背景を生成中 - お待ちください (Generating background locally - please wait)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 1", customBG, localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 3", customBG, localGenerate)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 譜面を解析中 (Analyzing chart)... ")
	levelData, err := pjsekaioverlay.FetchLevelData(chartSource, chart)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	var scoreMode string
	switch scoreModeInt {
	default:
		scoreMode = "default"
	case 2:
		scoreMode = "tournament"
	}
	if !isOptionSpecified {
		fmt.Print("\n採点モードを選択してください。(Choose scoring mode.)\n'1': デフォルト/Default\n'2': 大会モード/Tournament Mode (PERFECT = +3)\n> ")
		before, _ := rawmode.Enable()
		tmpScoreModeByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpScoreMode := string(tmpScoreModeByte)
		rawmode.Restore(before)
		switch tmpScoreMode {
		default:
			scoreMode = "default"
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpScoreMode))
			fmt.Println(color.GreenString("Score Mode: デフォルト/Default"))
		case "2":
			scoreMode = "tournament"
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpScoreMode))
			fmt.Println(color.GreenString("Score Mode: 大会/Tournament"))
		}
	}

	if !isOptionSpecified && scoreMode == "default" {
		fmt.Print("\n総合力を指定してください。 (Input your team power.)\n\n- 小数と科学的記数法が使える (Accepts decimals & scientific notation)\n- おすすめ (Recommended): 250000 - 300000\n- 例 (Example): 1234567; 1e+20; -300000\n> ")
		var tmpTeamPower string
		fmt.Scanln(&tmpTeamPower)
		if tmpTeamPower == "" {
			tmpTeamPower = "250000"
		}
		teamPower, err = strconv.ParseFloat(tmpTeamPower, 64)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		if teamPower >= math.Abs(math.Pow(2, 56)/10) && aviutlProcess == "aviutl.exe" {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.HiYellowString(tmpTeamPower))
			fmt.Println(color.HiYellowString("WARN: スコアは大きすぎると精度が落ちる可能性がある。Score may decrease precision if it's too large."))
		} else {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(tmpTeamPower))
		}
	}

	fmt.Print("- スコアを計算中 (Calculating score)... ")
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower, scoreMode, allFlick, listing)

	fmt.Println(color.GreenString("OK"))
	if !isOptionSpecified {
		fmt.Print("\n英語UIを使う？（部分的な対応）[y/n]\nUse English UI? (Partial support) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableENByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableEN := string(tmpEnableENByte)
		rawmode.Restore(before)
		if tmpEnableEN == "Y" || tmpEnableEN == "y" {
			enUI = true
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableEN))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			enUI = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableEN))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}

	executableDir := filepath.Dir(executablePath)
	assets := filepath.Join(executableDir, "assets")

	fmt.Print("- pedファイルを生成中 (Generating ped file)... ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, filepath.Join(formattedOutDir, "data.ped"), sonolus.LevelInfo{Rating: chart.Rating}, levelData, scoreMode, enUI, listing)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	var exoType = "exo"
	if aviutlProcess == "aviutl2.exe" {
		exoType = "alias(.object)"
	}

	fmt.Printf("- %sファイルを生成中 (Generating %s file)... ", exoType, exoType)

	var difficulty string
	difficultyStrings := []string{"EASY", "NORMAL", "HARD", "EXPERT", "MASTER", "APPEND", "ETERNAL"}

	for i := range chart.Tags {
		tags := checkSubstrings([]string{strings.ToUpper(chart.Tags[i].Title)}, difficultyStrings...)
		if tags != "" {
			difficulty = tags
			break
		}
	}

	if difficulty == "" {
		if title := checkSubstrings(strings.Fields(strings.ToUpper(chart.Title)), difficultyStrings...); title != "" {
			difficulty = title
		} else {
			difficulty = "APPEND"
		}
	}

	composerAndVocals := []string{chart.Artists, "-"}
	if separateAttempt := strings.Split(chart.Artists, " / "); chartSource.Id == "chart_cyanvas" && len(separateAttempt) <= 2 {
		composerAndVocals = separateAttempt
	}

	charter := []string{chart.Author, "-"}
	if charterTag := strings.Split(chart.Author, "#"); len(charterTag) <= 2 {
		charter = charterTag
	}

	description := []string{fmt.Sprintf("作詞：-    作曲：%s    編曲：-", composerAndVocals[0]), fmt.Sprintf("Vo：%s", composerAndVocals[1]), fmt.Sprintf("譜面ID：%s", strings.TrimPrefix(chartId, prefixTrim)), fmt.Sprintf("%s", charter[0])}
	descriptionv1 := []string{fmt.Sprintf("作詞：-    作曲：%s    編曲：-", composerAndVocals[0]), fmt.Sprintf("歌：%s    譜面制作：%s", composerAndVocals[1], charter[0])}
	extra := "【追加情報】"
	exFile := "tournament-mode.png"
	exFileOpacity := "100.0"

	if enUI {
		description = []string{fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -", composerAndVocals[0]), fmt.Sprintf("Vo：%s", composerAndVocals[1]), fmt.Sprintf("Chart ID: %s", strings.TrimPrefix(chartId, prefixTrim)), fmt.Sprintf("%s", charter[0])}
		descriptionv1 = []string{fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -", composerAndVocals[0]), fmt.Sprintf("Vocals: %s    Chart Design: %s", composerAndVocals[1], charter[0])}
		extra = "【Additional Info】"
		exFile = "tournament-mode-en.png"
	}
	if scoreMode == "tournament" {
		exFileOpacity = "0.0"
	}

	if aviutlProcess == "aviutl.exe" {
		err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, chart.Title, chart.Title, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
	} else {
		err = pjsekaioverlay.WriteAliasFiles(assets, formattedOutDir, chart.Title, chart.Title, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
	}

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	message := fmt.Sprintf("\n全ての処理が完了しました！READMEの規約を確認した上で、%sファイルを%sにインポートして下さい。\nExecution complete! Please import the %s file into %s after reviewing the README Terms of Use.", exoType, aviutlName, exoType, aviutlName)
	fmt.Println(color.GreenString(message))

	if !isOptionSpecified || !noExplorerAutoOpen {
		cmd := exec.Command(`explorer`, `/select,`, resultDir)
		cmd.Run()

		time.Sleep(2000 * time.Millisecond)
	}
}

func main() {
	isOptionSpecified := len(os.Args) > 1
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	origMain(isOptionSpecified)

	if !isOptionSpecified {
		fmt.Print(color.CyanString("\n- 何かキーを押すと終了します...\n- Press any key to exit..."))

		before, _ := rawmode.Enable()
		bufio.NewReader(os.Stdin).ReadByte()
		rawmode.Restore(before)
	}
}
