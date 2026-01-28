package main

import (
	"bufio"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/pjsekaioverlay"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
	"github.com/fatih/color"
	"github.com/google/go-github/v57/github"
	"github.com/srinathh/gokilo/rawmode"
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

//go:embed banlist.txt
var banUrl string

func BanList(name string) (bool, error) {
	resp, err := http.Get(banUrl)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	bodyStr := string(body)
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "4") || strings.HasPrefix(strings.TrimSpace(bodyStr), "5") {
		return false, fmt.Errorf("ban list error: %s", bodyStr)
	}

	banList := strings.Split(string(body), "\n")
	for _, bannedName := range banList {
		hashtagCount := strings.Count(name, "#")
		suffix := "#" + strings.Split(name, "#")[int(math.Max(0, float64(hashtagCount)-1))]

		if strings.TrimSpace(bannedName) == name {
			return true, nil
		} else if strings.HasSuffix(strings.TrimSpace(bannedName), suffix) {
			return true, nil
		} else if strings.EqualFold(strings.TrimSpace(bannedName), strings.TrimSuffix(name, suffix)) {
			return true, nil
		}
	}

	return false, nil
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
	flag.StringVar(&chartInstance, "instance", "", "サーバーインスタンスを指定します。(Specify the server instance.)")

	var customBG bool
	flag.BoolVar(&customBG, "custom-bg", false, "UntitledChartsでカスタム背景を使用する。(Use custom background in UntitledCharts.)")

	var scoreModeInt int
	flag.IntVar(&scoreModeInt, "score-mode", 1, "採点モードを指定します。(Specify scoring mode.)\n'1': デフォルト/Default\n'2': 大会モード/Tournament Mode (PERFECT = +3)\n'3': Sonolusスコア/Arcade Score (MAX: 1000000)")

	var teamPower float64
	flag.Float64Var(&teamPower, "power", 250000, "総合力を指定します。(Specify the team's power.)")

	var enUI bool
	flag.BoolVar(&enUI, "en-ui", false, "英語版を使う(イントロ + v3 UI) - Use English version (Intro + v3 UI)")

	var allFlick bool
	flag.BoolVar(&allFlick, "all-flick", false, "すべてのノーツをフリックとして扱います。(Treat all notes as flicks.)")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay-APPEND [オプション (Options)] [譜面ID (Chart ID)]")
		flag.PrintDefaults()
	}

	flag.Parse()

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

	if len(mapping) != 22 {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL:「default.ini」ファイルのデータに異常があります。「default.ini」ファイルを削除し、プログラムを再起動して再生成してください。\nAbnormal \"default.ini\" data. Please regenerate by deleting the \"default.ini\" file and reopen the program.\n- Mapping count: %v != 22", len(mapping))))
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
		"life":       mappingFloat64[5] >= 0 && mappingFloat64[5] <= 9999 && math.Mod(mappingFloat64[5], 1.0) == 0,
		"life_skill": mappingFloat64[6] == 0 || mappingFloat64[6] == 1,
		"overflow":   mappingFloat64[7] == 0 || mappingFloat64[7] == 1,
		"lead_zero":  mappingFloat64[8] == 0 || mappingFloat64[8] == 1,
		// Score
		"min_digit":   mappingFloat64[9] >= 1 && mappingFloat64[9] <= 99 && math.Mod(mappingFloat64[9], 1.0) == 0,
		"score_skill": mappingFloat64[10] >= 0 && mappingFloat64[10] <= 2 && math.Mod(mappingFloat64[10], 1.0) == 0,
		"score_speed": mappingFloat64[11] >= 0,
		"anim_score":  mappingFloat64[12] == 0 || mappingFloat64[12] == 1,
		"wds_anim":    mappingFloat64[13] == 0 || mappingFloat64[13] == 1,
		// Combo
		"ap":               mappingFloat64[14] == 0 || mappingFloat64[14] == 1,
		"tag":              mappingFloat64[15] == 0 || mappingFloat64[15] == 1,
		"last_digit":       mappingFloat64[16] >= 0 && math.Mod(mappingFloat64[16], 1.0) == 0,
		"combo_speed":      mappingFloat64[17] >= 0,
		"combo_burst":      mappingFloat64[18] == 0 || mappingFloat64[18] == 1,
		"achievement_rate": float64Pointer(mappingFloat64[19]) != nil,
		// Judgement
		"judge":       mappingFloat64[20] >= 1 && mappingFloat64[20] <= 10 && math.Mod(mappingFloat64[20], 1.0) == 0,
		"judge_speed": mappingFloat64[21] >= 0,
	}

	var mappingErr []string
	for i := range mapping {
		inRangeBool := inRange[mappingName[i]]
		if !inRangeBool {
			mappingErr = append(mappingErr, mappingName[i], fmt.Sprintf("%v", mapping[i]))
		}
	}

	if mappingErr != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:「default.ini」ファイルのデータに異常があります。「default.ini」ファイルを削除し、プログラムを再起動して再生成してください。\nAbnormal \"default.ini\" data. Please regenerate by deleting the \"default.ini\" file and reopen the program.\n- Mapping out of range: %s", mappingErr)))
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
		fmt.Print("譜面IDを接頭辞込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'sekai-rush-': Proseka Rush (sekairush.shop)\n'chcy-': Chart Cyanvas\n'ptlv-': Potato Leaves (ptlv.sevenc7c.com)\n'utsk-': Untitled Sekai (us.pim4n-net.com)\n'UnCh-': UntitledCharts (untitledcharts.com)\n'coconut-next-sekai-': Next SEKAI (coconut.sonolus.com/next-sekai)\n'lalo-': laoloser's server (sonolus.laoloser.com)\n'skyra-': osciris's server (Skyra)\n'sync-': Local Server (ScoreSync + ScoreSync Modern)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	// Instance section
	if chartInstance == "" && strings.HasPrefix(chartId, "chcy-") {
		fmt.Printf("\nChart Cyanvasインスタンスを選択してください。(Please choose Chart Cyanvas instance.)\n%s\n\n[インスタンス一覧 - List of instance(s)]\n'0': アーカイブ/Archive - cc.sevenc7c.com\n'1': 分岐サーバー/Offshoot server - chart-cyanvas.com\n> ", color.HiYellowString("(!) 別のインスタンスを持っていますか？URLドメインを入力してください。(Do you have a different instance? Input the URL domain.)"))
		var chartInput string
		fmt.Scanln(&chartInput)
		chartInput = strings.TrimPrefix(chartInput, "http://")
		chartInput = strings.TrimPrefix(chartInput, "https://")
		chartInstance = strings.Split(chartInput, "/")[0]
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartInput))
	} else if chartInstance == "" && strings.HasPrefix(chartId, "skyra-") {
		fmt.Printf("\nSkyraトークンを入力してください。(Please enter your Skyra token.)\n%s\n> ", color.HiYellowString("(!) トークンには「S-」という接頭辞が必要です。(Token must have a prefix of 'S-'.)"))
		var chartInput string
		fmt.Scanln(&chartInput)
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
		if chartSource.Status == 1 {
			fmt.Printf(color.RedString("FAIL: %sは対応されなくなりました。ご利用ありがとうございました。\n%s is no longer supported. Thank you for using the service.\n"), chartSource.Name, chartSource.Name)
			return
		}
		if chartSource.Status == 2 {
			fmt.Printf(color.HiYellowString("WARN: %sは現在開発中であり、正常に動作しない可能性があります。\n%s is currently in development and may not work.\n"), chartSource.Name, chartSource.Name)
		}
	}

	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())

	var chart sonolus.LevelInfo
	prefixTrim := checkSubstrings([]string{chartId}, "lalo-", "skyra-")
	chart, err = pjsekaioverlay.FetchChart(chartSource, strings.TrimPrefix(chartId, prefixTrim))

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	// Additional BG
	chartCCv1, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?c_background=v1")
	chartUNv3, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=v3")
	chartUNv1, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=v1")
	chartUNv1def, _ := pjsekaioverlay.FetchChart(chartSource, chartId+"?levelbg=default_or_v1")

	if chart.Engine.Version != 13 {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL (ver.%d): エンジンのバージョンが古い。pjsekai-overlay-APPENDを最新版に更新してください。\nUnsupported engine version. Please update pjsekai-overlay-APPEND to latest version.", chart.Engine.Version)))
		return
	}

	banList, err := BanList(chart.Author)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	} else if banList {
		fmt.Println(color.RedString("\nFAIL: 申し訳ありませんが、この譜面作者／組織はこのツールの使用が禁止されています。\nSorry, this charter/organization is banned from using this tool."))
		return
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

	// fmt.Print("- 音声のプレビューをダウンロード中 (Downloading preview audio)... ")
	// err = pjsekaioverlay.DownloadPreview(chartSource, chart, formattedOutDir)
	// if err != nil {
	// 	fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
	// 	return
	// }

	// fmt.Println(color.GreenString("OK"))

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

	if customBG {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		if chartSource.Id == "untitledcharts" {
			err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1def, formattedOutDir, chartId+"?levelbg=default_or_v1", "", customBG)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		} else {
			err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId+"/", "", customBG)
			if err != nil {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		}
	} else if chartSource.Id == "untitledcharts" {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv3, formattedOutDir, chartId+"?levelbg=v3", "", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1, formattedOutDir, chartId+"?levelbg=v1", "", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else if chartSource.Id == "chart_cyanvas" && chartSource.Name != "Chart Cyanvas Archive" {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartCCv1, formattedOutDir, chartId+"?c_background=v1", "", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else {
		fmt.Print("- ローカルで背景を生成中 - お待ちください (Generating background locally - please wait)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 1", customBG)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 3", customBG)
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
	case 3:
		scoreMode = "sonolus"
	}
	if !isOptionSpecified {
		fmt.Print("\n採点モードを選択してください。(Choose scoring mode.)\n'1': デフォルト/Default\n'2': 大会モード/Tournament Mode (PERFECT = +3)\n'3': Sonolusスコア/Arcade Score (MAX: 1000000)\n> ")
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
		case "3":
			scoreMode = "sonolus"
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpScoreMode))
			fmt.Println(color.GreenString("Score Mode: Sonolusスコア/Arcade Score"))
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
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower, scoreMode, allFlick)

	fmt.Println(color.GreenString("OK"))
	if !isOptionSpecified {
		fmt.Print("\n英語UIを使う？（イントロ + v3 UI）[y/n]\nUse English UI? (Intro + v3 UI) [y/n]\n> ")
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

	fmt.Print("\n- pedファイルを生成中 (Generating ped file)... ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, filepath.Join(formattedOutDir, "data.ped"), sonolus.LevelInfo{Rating: chart.Rating}, levelData, scoreMode, enUI)

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

	description := []string{fmt.Sprintf("作詞：-    作曲：%s    編曲：-", composerAndVocals[0]), fmt.Sprintf("Vo：%s    譜面制作：%s", composerAndVocals[1], charter[0])}
	descriptionv1 := []string{fmt.Sprintf("作詞：-    作曲：%s    編曲：-", composerAndVocals[0]), fmt.Sprintf("歌：%s    譜面制作：%s", composerAndVocals[1], charter[0])}
	extra := "【追加情報】"
	exFile := "tournament-mode.png"
	exFileOpacity := "100.0"

	if enUI {
		description = []string{fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -", composerAndVocals[0]), fmt.Sprintf("Vo: %s    Chart Design: %s", composerAndVocals[1], charter[0])}
		descriptionv1 = []string{fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -", composerAndVocals[0]), fmt.Sprintf("Vocals: %s    Chart Design: %s", composerAndVocals[1], charter[0])}
		extra = "【Additional Info】"
		exFile = "tournament-mode-en.png"
	}
	if scoreMode == "tournament" {
		exFileOpacity = "0.0"
	}

	if aviutlProcess == "aviutl.exe" {
		err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, chart.Title, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
	} else {
		err = pjsekaioverlay.WriteAliasFiles(assets, formattedOutDir, chart.Title, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, mappingStr)
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
