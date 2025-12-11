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

func BanList(name string) bool {
	resp, err := http.Get(banUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
	}

	bodyStr := string(body)
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "4") || strings.HasPrefix(strings.TrimSpace(bodyStr), "5") {
		panic(bodyStr)
	}

	banList := strings.Split(string(body), "\n")
	for _, bannedName := range banList {
		if strings.TrimSpace(bannedName) == name {
			return true
		} else if strings.HasSuffix(strings.TrimSpace(bannedName), "#"+strings.Split(name, "#")[1]) {
			return true
		} else if strings.EqualFold(strings.TrimSpace(bannedName), strings.Split(name, "#")[0]) {
			return true
		}
	}

	return false
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

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./dist/_chartId_", "出力先ディレクトリを指定します。_chartId_ は譜面IDに置き換えられます。\nEnter the output path. _chartId_ will be replaced with the chart ID.")

	var chartInstance string
	flag.StringVar(&chartInstance, "instance", "", "サーバーインスタンスを指定します。(Specify the server instance.)")

	var customBG bool
	flag.BoolVar(&customBG, "custom-bg", false, "UntitledChartsでカスタム背景を使用する。(Use custom background in UntitledCharts.)")

	var scoreMode string
	flag.StringVar(&scoreMode, "score-mode", "default", "採点モードを指定します。(Specify scoring mode.)")

	var teamPower float64
	flag.Float64Var(&teamPower, "power", 250000, "総合力を指定します。(Specify the team's power.)")

	var enUI bool
	flag.BoolVar(&enUI, "en-ui", false, "英語版を使う(イントロ + v3 UI) - Use English version (Intro + v3 UI)")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay-APPEND [譜面ID (Chart ID)] [オプション (Options)]")
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

	locale := func() string {
		cmd := exec.Command("powershell", "-Command", "Get-WinSystemLocale | Select-Object -ExpandProperty Name")
		output, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(output))
	}
	if locale() != "ja-JP" {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: お使いのシステムロケールが「日本語（日本）」に設定されていません。変更方法についてはREADMEを参照してください。\nYour system locale is not set to \"Japanese (Japan)\". Refer to README for how to change it.\n- System locale: %v", locale())))
		return
	}

	mappingName, mapping := pjsekaioverlay.SetOverlayDefault()

	if len(mapping) != 22 {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL:「default.ini」ファイルのデータに異常があります。「default.ini」ファイルを削除し、プログラムを再起動して再生成してください。\nAbnormal \"default.ini\" data. Please regenerate by deleting the \"default.ini\" file and reopen the program.\n- Mapping count: %v != 22", len(mapping))))
		return
	}

	var inRange = map[string]bool{
		// Root
		"offset":      mapping[0] >= -99999.99 && mapping[0] <= 99999.99,
		"cache":       mapping[1] == 0 || mapping[1] == 1,
		"font_type":   mapping[2] == 0 || mapping[2] == 1,
		"text_lang":   mapping[3] == 0 || mapping[3] == 1,
		"watermark":   mapping[4] == 0 || mapping[4] == 1,
		"detail_stat": mapping[5] == 0 || mapping[5] == 1,
		// Life
		"life":       mapping[6] >= 0 && mapping[6] <= 9999 && math.Mod(mapping[6], 1.0) == 0,
		"life_skill": mapping[7] == 0 || mapping[7] == 1,
		"overflow":   mapping[8] == 0 || mapping[8] == 1,
		"lead_zero":  mapping[9] == 0 || mapping[9] == 1,
		// Score
		"min_digit":   mapping[10] >= 1 && mapping[10] <= 99 && math.Mod(mapping[10], 1.0) == 0,
		"score_skill": mapping[11] >= 0 && mapping[11] <= 2 && math.Mod(mapping[11], 1.0) == 0,
		"score_speed": mapping[12] >= 0,
		"anim_score":  mapping[13] == 0 || mapping[13] == 1,
		"wds_anim":    mapping[14] == 0 || mapping[14] == 1,
		// Combo
		"ap":          mapping[15] == 0 || mapping[15] == 1,
		"tag":         mapping[16] == 0 || mapping[16] == 1,
		"last_digit":  mapping[17] >= 0 && math.Mod(mapping[17], 1.0) == 0,
		"combo_speed": mapping[18] >= 0,
		"combo_burst": mapping[19] == 0 || mapping[19] == 1,
		// Judgement
		"judge":       mapping[20] >= 1 && mapping[20] <= 10 && math.Mod(mapping[20], 1.0) == 0,
		"judge_speed": mapping[21] >= 0,
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
		fmt.Print("譜面IDを接頭辞込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'chcy-': Chart Cyanvas\n'ptlv-': Potato Leaves (ptlv.sevenc7c.com)\n'utsk-': Untitled Sekai (us.pim4n-net.com)\n'UnCh-': UntitledCharts (untitledcharts.com)\n'coconut-next-sekai-': Next SEKAI (coconut.sonolus.com/next-sekai)\n'lalo-': laoloser's server (sonolus.laoloser.com)\n'sync-': Local Server (ScoreSync)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	if chartInstance == "" && strings.HasPrefix(chartId, "chcy-") {
		fmt.Printf("\nChart Cyanvasインスタンスを選択してください。(Please choose Chart Cyanvas instance.)\n%s\n\n[インスタンス一覧 - List of instance(s)]\n'0': アーカイブ/Archive - cc.sevenc7c.com\n'1': 分岐サーバー/Offshoot server - chart-cyanvas.com\n> ", color.HiYellowString("(!) 別のインスタンスを持っていますか？URLドメインを入力してください。(Do you have a different instance? Input the URL domain.)"))
		var chartInput string
		fmt.Scanln(&chartInput)
		chartInput = strings.TrimPrefix(chartInput, "http://")
		chartInput = strings.TrimPrefix(chartInput, "https://")
		chartInstance = strings.Split(chartInput, "/")[0]
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartInput))
	}

	var chartSource pjsekaioverlay.Source
	var err error
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

	if !isOptionSpecified && chartSource.Id == "untitledcharts" {
		message := fmt.Sprintf("NOTE: %sでの録画には別の動画の作り方が必要です。詳細はREADMEの「動画の作り方」をご確認ください。\nRecording in %s require a different method for creating videos. Please check \"Video Guide\" in README for details.\n\n- 何かキーを押すと続行します...\n- Press any key to continue...", chartSource.Name, chartSource.Name)
		fmt.Println(color.CyanString(message))

		before, _ := rawmode.Enable()
		bufio.NewReader(os.Stdin).ReadByte()
		rawmode.Restore(before)
	}

	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())

	var chart sonolus.LevelInfo
	if strings.HasPrefix(chartId, "lalo-") {
		chart, err = pjsekaioverlay.FetchChart(chartSource, chartId[5:])
	} else {
		chart, err = pjsekaioverlay.FetchChart(chartSource, chartId)
	}

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

	if BanList(chart.Author) {
		fmt.Println(color.RedString("\nFAIL: 申し訳ありませんが、この譜面作者／組織はこのツールの使用が禁止されています。\n組織のリーダーである場合、および／または禁止措置の異議申し立てを希望する場合は、非公開で連絡してください。この状況を公表すると、上訴が無効となり、その状態が永久に続くことにご留意ください。\n\nSorry, this charter/organization is banned from using this tool.\nIf you're the leader of the organization and/or would like to request a ban appeal, you must contact me privately. Please note that publicizing this situation will nullify your appeal indefinitely."))
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

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.ReplaceAll(outDir, "_chartId_", chartId))
	resultDir := filepath.Dir(formattedOutDir) + "\\" + chartId

	isAdminPerm := func(path string) bool {
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
	if isAdminPerm(formattedOutDir) {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL: ディレクトリには管理者権限が必要です。pjsekai-overlay-APPENDを「C:\\」または別の場所に移動してください。\nYour directory requires administrative permissions. Please move pjsekai-overlay-APPEND to \"C:\\\" or somewhere else.\n\n出力先ディレクトリ (Output path): %s", resultDir)))
		return
	}

	isASCII := func(s string) bool {
		for i := 0; i < len(s); i++ {
			if s[i] > 127 {
				return false
			}
		}
		return true
	}
	if !isASCII(formattedOutDir) {
		fmt.Println(color.RedString(fmt.Sprintf("\nFAIL: ディレクトリに非ASCII文字が含まれています。pjsekai-overlay-APPENDを「C:\\」または別の場所に移動してください。\nYour directory contains non-ASCII characters. Please move pjsekai-overlay-APPEND to \"C:\\\" or somewhere else.\n\n出力先ディレクトリ (Output path): %s", resultDir)))
		return
	}

	fmt.Println(color.GreenString("OK"))
	fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(resultDir))

	fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
	err = pjsekaioverlay.DownloadCover(chartSource, chart, formattedOutDir)
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

	if !isOptionSpecified && chartSource.Id == "untitledcharts" {
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

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1def, formattedOutDir, chartId+"?levelbg=default_or_v1", "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else if chartSource.Id == "untitledcharts" {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv3, formattedOutDir, chartId+"?levelbg=v3", "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartUNv1, formattedOutDir, chartId+"?levelbg=v1", "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else if chartSource.Id == "chart_cyanvas" && chartSource.Name != "Chart Cyanvas Archive" {
		fmt.Print("- 背景をダウンロード中 (Downloading background)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chartCCv1, formattedOutDir, chartId+"?c_background=v1", "")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}
	} else {
		fmt.Print("- ローカルで背景を生成中 - しばらく時間がかかります (Generating background locally - will take a while)... ")

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 1")
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
			return
		}

		err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId, "-v 3")
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
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower, scoreMode)

	fmt.Println(color.GreenString("OK"))
	if !isOptionSpecified {
		if aviutlProcess == "aviutl.exe" {
			fmt.Print("\n英語UIを使う？（設定@pjsekai-overlayテキスト + イントロ + v3 UI）[y/n]\nUse English UI? (Root@pjsekai-overlay text + Intro + v3 UI) [y/n]\n> ")
		} else {
			fmt.Print("\n英語UIを使う？（イントロ + v3 UI）[y/n]\nUse English UI? (Intro + v3 UI) [y/n]\n> ")
		}
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

	cmd := exec.Command(`explorer`, `/select,`, resultDir)
	cmd.Run()

	time.Sleep(2000 * time.Millisecond)
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
