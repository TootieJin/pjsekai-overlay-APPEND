package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/pjsekaioverlay"
	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
	"github.com/fatih/color"
	"github.com/google/go-github/v57/github"
	"github.com/srinathh/gokilo/rawmode"
	"golang.org/x/sys/windows"
)

func checkUpdate() []string {
	githubClient := github.NewClient(nil)
	release, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), "TootieJin", "pjsekai-overlay-APPEND")
	if err != nil {
		return nil
	}

	latestVersion := strings.TrimPrefix(release.GetTagName(), "v")
	if latestVersion == pjsekaioverlay.Version || pjsekaioverlay.Version == "0.0.0" {
		return nil
	}
	return []string{latestVersion, release.GetHTMLURL()}
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

func origMain(isOptionSpecified bool) {
	Title()

	var skipAviutlInstall bool
	flag.BoolVar(&skipAviutlInstall, "no-aviutl-install", false, "AviUtlオブジェクトのインストールをスキップします。(AviUtl object installation is skipped.)")

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./dist/_chartId_", "出力先ディレクトリを指定します。_chartId_ は譜面IDに置き換えられます。\nEnter the output path. _chartId_ will be replaced with the chart ID.")

	var exScore bool
	flag.BoolVar(&exScore, "ex-score", false, "大会モードを有効にします。 (Enable Tournament Mode.)")

	var teamPower float64
	flag.Float64Var(&teamPower, "team-power", 250000, "総合力を指定します。(Enter the team's power.)")

	var enUI bool
	flag.BoolVar(&enUI, "en-ui", false, "英語版を使う(イントロ + v3 UI) - Use English version (Intro + v3 UI)")

	var apCombo bool
	flag.BoolVar(&apCombo, "ap-combo", true, "コンボのAP表示を有効にします。(Enable AP display for combo.)")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay [譜面ID] [オプション]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if checkUpdate() != nil {
		fmt.Printf(color.HiCyanString("新しいバージョンがリリースされています\nNew version released: v%s -> v%s\n"), pjsekaioverlay.Version, checkUpdate()[0])
		fmt.Printf(color.HiCyanString("ダウンロード (Download Here) -> %s\n"), checkUpdate()[1])
		fmt.Println(color.RedString("\nFAIL: pjsekai-overlay-APPENDを最新バージョンに更新してください。\nUpdate pjsekai-overlay-APPEND to the latest version."))
		return
	}

	if !skipAviutlInstall {
		success := pjsekaioverlay.TryInstallObject()
		if success {
			fmt.Println(color.GreenString("AviUtlオブジェクトのインストールに成功しました。(AviUtl object successfully installed.)\n"))
		}
	}

	Tips()

	var chartId string
	if flag.Arg(0) != "" {
		chartId = flag.Arg(0)
		fmt.Printf("譜面ID (Chart ID): %s\n", color.GreenString(chartId))
	} else {
		fmt.Print("譜面IDを接頭辞込みで入力して下さい。\nEnter the chart ID including the prefix.\n\n'chcy-': Chart Cyanvas\n'ptlv-': Potato Leaves (ptlv.sevenc7c.com)\n'utsk-': Untitled Sekai (us.pim4n-net.com)\n'UnCh-': UntitledCharts (untitledcharts.com)\n'coconut-next-sekai-': Next SEKAI (coconut.sonolus.com/next-sekai)\n'sync-': Local Server (ScoreSync)\n> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(chartId))
	}

	var chartInstance []string
	if strings.HasPrefix(chartId, "chcy-") {
		fmt.Printf("\nChart Cyanvasインスタンスを選択してください。(Please choose Chart Cyanvas instance.)\n%s\n\n[インスタンス一覧 - List of instance(s)]\n'0': アーカイブ/Archive - cc.sevenc7c.com\n'1': 分岐サーバー/Offshoot server - chart-cyanvas.com\n> ", color.YellowString("(!) 別のインスタンスを持っていますか？URLドメインを入力してください。(Do you have a different instance? Input the URL domain.)"))
		var chartInput string
		fmt.Scanln(&chartInput)
		chartInput = strings.TrimPrefix(chartInput, "http://")
		chartInput = strings.TrimPrefix(chartInput, "https://")
		chartInstance = strings.Split(chartInput, "/")
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
			fmt.Print("ローカルサーバーの譜面タイトルを入力してください。(Enter chart title for the local server.)\n> ")
			fmt.Scanln(&chartId)
		}
	} else {
		chartSource, err = pjsekaioverlay.DetectChartSource(chartId, chartInstance)
		if err != nil {
			fmt.Println(color.RedString("FAIL: 譜面が見つかりません。接頭辞も込め、正しい譜面IDを入力して下さい。\nChart not found. Please enter the correct chart ID including the prefix."))
			return
		}
		if chartSource.Status == 1 {
			fmt.Printf(color.RedString("FAIL: %sはサポートされなくなりました。ご利用ありがとうございました。\n%s is no longer supported. Thank you for using the service.\n"), chartSource.Name, chartSource.Name)
			return
		}
		if chartSource.Status == 2 {
			fmt.Printf(color.HiYellowString("WARN: %sは現在開発中であり、正常に動作しない可能性があります。\n%s is currently in development and may not work.\n"), chartSource.Name, chartSource.Name)
		}
	}

	fmt.Printf("- 譜面を取得中 (Getting chart): %s%s%s ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())
	chart, err := pjsekaioverlay.FetchChart(chartSource, chartId)
	chartv1, errv1 := pjsekaioverlay.FetchChart(chartSource, chartId+"?c_background=v1")

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	if errv1 != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", errv1.Error())))
		return
	}

	if chart.Engine.Version != 13 {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL (ver.%d):エンジンのバージョンが古い。pjsekai-overlay-APPENDを最新版に更新してください。\nUnsupported engine version. Please update pjsekai-overlay-APPEND to latest version.", chart.Engine.Version)))
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

	fmt.Println(color.GreenString("OK"))

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.Replace(outDir, "_chartId_", chartId, -1))
	fmt.Printf("- 出力先ディレクトリ (Output path): %s\n", color.CyanString(filepath.Dir(formattedOutDir)+"\\"+chartId))

	fmt.Print("- ジャケットをダウンロード中 (Downloading jacket)... ")
	err = pjsekaioverlay.DownloadCover(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- 背景をダウンロード中 (Downloading background)... ")
	err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir, chartId)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}
	err = pjsekaioverlay.DownloadBackground(chartSource, chartv1, formattedOutDir, chartId+"?c_background=v1")
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
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
		fmt.Print("\n大会モードを有効にするか？ (PERFECT = +3点)\nEnable Tournament Mode? (PERFECT = +3pts) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableEXScoreByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableEXScore := string(tmpEnableEXScoreByte)
		rawmode.Restore(before)
		if tmpEnableEXScore == "Y" || tmpEnableEXScore == "y" {
			exScore = true
			teamPower = 0.0
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableEXScore))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			exScore = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableEXScore))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}

	if !isOptionSpecified && !exScore {
		fmt.Print("\n総合力を指定してください。 (Input your team power.)\n\n- 小数と科学的記数法が使える (Accepts decimals & scientific notation)\n- おすすめ (Recommended): 250000 - 300000\n- 制限 (Limit): ???\n> ")
		var tmpTeamPower string
		fmt.Scanln(&tmpTeamPower)
		teamPower, err = strconv.ParseFloat(tmpTeamPower, 64)
		if err != nil {
			if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
				fmt.Println(color.RedString("FAIL: あなたのPCがその総合力で計算できないのは残念だ。説明書を読んで再実行してください。\nToo bad your PC can't calculate with that team power. Read the instructions and rerun it."))
				return
			} else {
				fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
				return
			}
		}
		if teamPower >= math.Abs(1e+33) {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.YellowString(tmpTeamPower))
			fmt.Println(color.HiYellowString("WARN: スコアは大きすぎると精度が落ちる可能性がある。Score may decrease precision if it's too large.\n"))
		} else {
			fmt.Printf("\033[A\033[2K\r> %s\n", color.GreenString(tmpTeamPower))
		}
	}

	fmt.Print("- スコアを計算中 (Calculating score)... ")
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower, exScore)

	fmt.Println(color.GreenString("OK"))
	if !isOptionSpecified {
		fmt.Print("\n英語版を使う？(イントロ + v3 UI) - Use English version? (Intro + v3 UI) [y/n]\n> ")
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

	if !isOptionSpecified {
		fmt.Print("\nコンボのAP表示を有効にしますか？(これは後でAviUtlで変更できます)\nEnable AP indicator for combo? (You can change this later in AviUtl) [y/n]\n> ")
		before, _ := rawmode.Enable()
		tmpEnableComboApByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableComboAp := string(tmpEnableComboApByte)
		rawmode.Restore(before)

		if tmpEnableComboAp == "Y" || tmpEnableComboAp == "y" || tmpEnableComboAp == "" {
			apCombo = true
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.GreenString(tmpEnableComboAp))
			fmt.Println(color.GreenString("TOGGLE: ON"))
		} else {
			apCombo = false
			fmt.Printf("\n\033[A\033[2K\r> %s\n", color.RedString(tmpEnableComboAp))
			fmt.Println(color.RedString("TOGGLE: OFF"))
		}
	}
	executableDir := filepath.Dir(executablePath)
	assets := filepath.Join(executableDir, "assets")

	fmt.Print("\n- pedファイルを生成中 (Generating ped file)... ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, filepath.Join(formattedOutDir, "data.ped"), sonolus.LevelInfo{Rating: chart.Rating}, levelData, exScore, enUI)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Print("- exoファイルを生成中 (Generating exo file)... ")

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
	if charterTag := strings.Split(chart.Author, "#"); (chartSource.Id == "chart_cyanvas" || chartSource.Id == "untitled_sekai" || chartSource.Id == "untitledcharts") && len(charterTag) <= 2 {
		charter = charterTag
	}

	description := fmt.Sprintf("作詞：-    作曲：%s    編曲：-\r\nVo：%s    譜面作成：%s", composerAndVocals[0], composerAndVocals[1], charter[0])
	descriptionv1 := fmt.Sprintf("作詞：-   作曲：%s   編曲：-\r\n歌：%s   譜面作成：%s", composerAndVocals[0], composerAndVocals[1], charter[0])
	extra := "[追加情報]"
	exFile := "tournament-mode.png"
	exFileOpacity := "100.0"
	ap := "0.00"

	if enUI {
		description = fmt.Sprintf("Lyrics: -    Music: %s    Arrangement: -\r\nVo: %s    Chart Design: %s", composerAndVocals[0], composerAndVocals[1], charter[0])
		descriptionv1 = fmt.Sprintf("Lyrics: -   Music: %s   Arrangement: -\r\n歌：%s   Chart Design: %s", composerAndVocals[0], composerAndVocals[1], charter[0])
		extra = "[Additional Info]"
		exFile = "tournament-mode-en.png"
	}
	if exScore {
		exFileOpacity = "0.0"
	}
	if apCombo {
		ap = "1.00"
	}

	err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir, chart.Title, description, descriptionv1, difficulty, extra, exFile, exFileOpacity, ap)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("FAIL: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("OK"))

	fmt.Println(color.GreenString("\n全ての処理が完了しました。READMEの規約を確認した上で、exoファイルをAviUtlにインポートして下さい。\nExecution complete! Please import the exo file into AviUtl after reviewing the README Terms of Use."))
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
