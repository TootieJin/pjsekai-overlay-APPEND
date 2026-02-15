package pjsekaioverlay

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
)

type Source struct {
	Id     string
	Name   string
	Color  int
	Host   string
	Status int
}

func DetectLocalChartSource() (Source, error) {
	// ScoreSync（3939ポート）に接続を試行
	resp, err := http.Get("http://localhost:3939/")
	if err != nil {
		return Source{}, errors.New("ScoreSyncに接続できませんでした。サーバーが起動していることを確認してください。(Could not connect to ScoreSync. Please make sure the server is up and running.)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Source{}, errors.New("ScoreSyncからの応答が正常ではありません。(Invalid response from ScoreSync.)")
	}

	source := Source{
		Id:     "local_server",
		Name:   "Local Server (ScoreSync + ScoreSync Modern)",
		Color:  0xa9cd6a,
		Host:   "localhost:3939",
		Status: 0,
	}
	return source, nil
}

func FetchChart(source Source, chartId string) (sonolus.LevelInfo, error) {
	var url string
	if source.Id == "local_server" {
		// ローカルサーバーの場合はchartIdをタイトルとして使用
		url = "http://" + source.Host + "/sonolus/levels/" + chartId
	} else {
		url = "https://" + source.Host + "/sonolus/levels/" + chartId
	}

	resp, err := http.Get(url)

	if err != nil {
		if source.Id == "local_server" {
			return sonolus.LevelInfo{}, errors.New("ScoreSyncに接続できませんでした。(Could not connect to ScoreSync.)")
		}
		return sonolus.LevelInfo{}, errors.New("サーバーに接続できませんでした。(Could not connect to server.)")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if source.Id == "local_server" {
			return sonolus.LevelInfo{}, errors.New("指定されたタイトルの譜面が見つかりませんでした。(Chart not found.)")
		}
		return sonolus.LevelInfo{}, errors.New("譜面が見つかりませんでした。(Unable to search chart.)")
	}

	var chart sonolus.InfoResponse[sonolus.LevelInfo]
	json.NewDecoder(resp.Body).Decode(&chart)

	return chart.Item, nil
}

func DetectChartSource(chartId string, chartInstance string) (Source, error) {
	var source Source
	if strings.HasPrefix(chartId, "sekai-rush-") {
		source = Source{
			Id:     "proseka_rush",
			Name:   "Proseka Rush",
			Color:  0x02cbbd,
			Host:   "sekairush.shop",
			Status: 0,
		}
	} else if strings.HasPrefix(chartId, "chcy-") {
		switch chartInstance {
		case "0":
			source = Source{
				Id:     "chart_cyanvas",
				Name:   "Chart Cyanvas Archive",
				Color:  0x83ccd2,
				Host:   "cc.sevenc7c.com",
				Status: 0,
			}
		case "1":
			source = Source{
				Id:     "chart_cyanvas",
				Name:   "Chart Cyanvas Offshoot Server",
				Color:  0x83ccd2,
				Host:   "chart-cyanvas.com",
				Status: 0,
			}
		default:
			source = Source{
				Id:     "chart_cyanvas",
				Name:   "Chart Cyanvas (" + chartInstance + ")",
				Color:  0x83ccd2,
				Host:   chartInstance,
				Status: 0,
			}
		}
	} else if strings.HasPrefix(chartId, "ptlv-") {
		source = Source{
			Id:     "potato_leaves",
			Name:   "Potato Leaves",
			Color:  0x88cb7f,
			Host:   "ptlv.sevenc7c.com",
			Status: 0,
		}
	} else if strings.HasPrefix(chartId, "utsk-") {
		source = Source{
			Id:     "untitled_sekai",
			Name:   "Untitled Sekai",
			Color:  0x6a6a6a,
			Host:   "us.pim4n-net.com",
			Status: 1,
		}
	} else if strings.HasPrefix(chartId, "UnCh-") {
		source = Source{
			Id:     "untitledcharts",
			Name:   "UntitledCharts",
			Color:  0x7765da,
			Host:   "untitledcharts.com",
			Status: 0,
		}
	} else if strings.HasPrefix(chartId, "coconut-next-sekai-") {
		source = Source{
			Id:     "next_sekai",
			Name:   "Next SEKAI",
			Color:  0x5b5c7c,
			Host:   "coconut.sonolus.com/next-sekai",
			Status: 0,
		}
	} else if strings.HasPrefix(chartId, "lalo-") {
		source = Source{
			Id:     "laoloser",
			Name:   "laoloser's server",
			Color:  0xccd1df,
			Host:   "sonolus.laoloser.com",
			Status: 1,
		}
	} else if strings.HasPrefix(chartId, "skyra-") {
		switch chartInstance {
		default:
			source = Source{
				Id:     "skyra",
				Name:   "Skyra",
				Color:  0x9c4edc,
				Host:   "skyra.plumnet.live",
				Status: 1,
			}
		}
	}
	if source.Id == "" {
		return Source{
			Id:     chartId,
			Name:   "",
			Color:  0,
			Host:   "",
			Status: 0,
		}, errors.New("unknown chart source")
	}
	return source, nil
}

func FetchLevelData(source Source, level sonolus.LevelInfo) (sonolus.LevelData, error) {
	var url string
	var err error

	if source.Id == "local_server" {
		url, err = sonolus.JoinUrl("http://"+source.Host, level.Data.Url)
	} else {
		url, err = sonolus.JoinUrl("https://"+source.Host, level.Data.Url)
	}

	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return sonolus.LevelData{}, fmt.Errorf("譜面データが見つかりませんでした。(No chart data found.) [%d]", resp.StatusCode)
	}

	var data sonolus.LevelData
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("譜面データの読み込みに失敗しました。(Loading chart data failed.) [%s]", err)
	}

	err = json.NewDecoder(gzipReader).Decode(&data)

	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("譜面データの読み込みに失敗しました。(Loading chart data failed.) [%s]", err)
	}

	return data, nil
}

func DownloadJacket(source Source, level sonolus.LevelInfo, destPath string) error {
	var url string
	var err error

	if source.Id == "local_server" {
		url, err = sonolus.JoinUrl("http://"+source.Host, level.Cover.Url)
	} else {
		url, err = sonolus.JoinUrl("https://"+source.Host, level.Cover.Url)
	}

	if err != nil {
		return fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%d]", resp.StatusCode)
	}

	os.MkdirAll(destPath, 0755)
	imageData, _, err := image.Decode(resp.Body)

	if err != nil {
		return fmt.Errorf("ジャケットの読み込みに失敗しました。(Loading jacket failed.) [%s]", err)
	}

	newImage := image.NewRGBA(image.Rect(0, 0, 512, 512))

	draw.ApproxBiLinear.Scale(newImage, newImage.Bounds(), imageData, imageData.Bounds(), draw.Over, nil)

	file, err := os.Create(path.Join(destPath, "cover.png"))

	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
	}

	defer file.Close()

	err = png.Encode(file, newImage)

	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
	}

	return nil
}

func DownloadPreview(source Source, level sonolus.LevelInfo, destPath string) error {
	var previewUrl string
	var err error

	if source.Id == "local_server" {
		previewUrl, err = sonolus.JoinUrl("http://"+source.Host, level.Preview.Url)
	} else {
		previewUrl, err = sonolus.JoinUrl("https://"+source.Host, level.Preview.Url)
	}

	if err != nil {
		return fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
	}

	resp, err := http.Get(previewUrl)
	if err != nil {
		return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("音声が見つかりませんでした。(Audio not found.) [%d]", resp.StatusCode)
	}

	var file *os.File
	file, err = os.Create(path.Join(destPath, "preview.mp3"))

	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
	}

	defer file.Close()

	if file != nil {
		if _, err := io.Copy(file, resp.Body); err != nil {
			return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
		}
	}
	return nil
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	return nil
}

func DownloadBackground(source Source, level sonolus.LevelInfo, destPath string, chartId string, arg string, customBG bool) error {
	if source.Id == "proseka_rush" || source.Name == "Chart Cyanvas Archive" || source.Id == "potato_leaves" || source.Id == "local_server" || source.Id == "next_sekai" || (source.Id == "skyra" && !customBG) {
		coverPath := path.Join(destPath, "cover.png")
		if _, err := os.Stat(coverPath); os.IsNotExist(err) {
			return fmt.Errorf("ジャケット画像が見つかりません。先にジャケット画像をダウンロードしてください。(Jacket image not found. Download jacket image first.)")
		}

		executablePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("実行ファイルパスの取得に失敗しました。(Failed to obtain execuable path.) [%s]", err)
		}

		backgroundGenPath := path.Join(path.Dir(executablePath), "dependencies", "pjsekai-background-gen.exe")

		var outputPath string
		if arg == "-v 1" {
			outputPath = path.Join(destPath, "background-v1.png")
		} else {
			outputPath = path.Join(destPath, "background.png")
		}

		// pjsekai-background-gen.exeが存在しない場合はダウンロード
		if _, err := os.Stat(backgroundGenPath); os.IsNotExist(err) {
			err = DownloadBackgroundGenerator(backgroundGenPath)
			if err != nil {
				return fmt.Errorf("背景生成ツールのダウンロードに失敗しました。(Failed to download background generator.) [%s]", err)
			}
		}

		absBackgroundGenPath, err := filepath.Abs(backgroundGenPath)
		if err != nil {
			return fmt.Errorf("背景生成ツールのパス解決に失敗しました。(Background generator failed to resolve path.) [%s]", err)
		}

		absCoverPath, err := filepath.Abs(coverPath)
		if err != nil {
			return fmt.Errorf("ジャケット画像のパス解決に失敗しました。(Failed to resolve path for jacket image.) [%s]", err)
		}

		cmd := exec.Command(absBackgroundGenPath, absCoverPath, arg)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("背景生成に失敗しました。(Failed to generate background.) [%s] \n出力/Output: \n%s\n コマンド/Command: %s", err, string(output), cmd)
		}

		generatedBackgroundPath := path.Join(destPath, "cover.output.png")

		if _, err := os.Stat(generatedBackgroundPath); os.IsNotExist(err) {
			if _, err := os.Stat(absCoverPath); err == nil {
				err = CopyFile(absCoverPath, outputPath)
				if err != nil {
					return fmt.Errorf("背景ファイルのコピーに失敗しました。(Failed to copy background file.)[%s]", err)
				}
				return nil
			}
			return fmt.Errorf("背景ファイルが生成されませんでした (Failed to generate background): %s", generatedBackgroundPath)
		}

		// cover.output.pngをbackground.pngにリネーム
		err = os.Rename(generatedBackgroundPath, outputPath)
		if err != nil {
			return fmt.Errorf("背景ファイルのリネームに失敗しました。(Failed to rename background file.) [%s]", err)
		}
	} else {
		var backgroundUrl string
		var err error

		if level.UseBackground.UseDefault {
			backgroundUrl, err = sonolus.JoinUrl("https://"+source.Host, level.Engine.Background.Image.Url)
		} else {
			backgroundUrl, err = sonolus.JoinUrl("https://"+source.Host, level.UseBackground.Item.Image.Url)
		}

		if err != nil {
			return fmt.Errorf("URLの解析に失敗しました。(URL parsing failed.) [%s]", err)
		}

		resp, err := http.Get(backgroundUrl)
		if err != nil {
			return fmt.Errorf("サーバーに接続できませんでした。(Could not connect to server.) [%s]", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("背景が見つかりませんでした。(Background not found.) [%d]", resp.StatusCode)
		}

		var file *os.File
		var filev1 *os.File

		if strings.Contains(chartId, "?c_background=v1") || strings.Contains(chartId, "?levelbg=default_or_v1") || strings.Contains(chartId, "?levelbg=v1") || strings.HasSuffix(chartId, "/") { // v1 BG (or custom)
			filev1, err = os.Create(path.Join(destPath, "background-v1.png"))
			file = nil
		} else {
			filev1 = nil
			file, err = os.Create(path.Join(destPath, "background.png"))
		}

		if err != nil {
			return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
		}

		defer file.Close()
		defer filev1.Close()

		if file != nil {
			if _, err := io.Copy(file, resp.Body); err != nil {
				return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
			}
		}
		if filev1 != nil {
			if _, err := io.Copy(filev1, resp.Body); err != nil {
				return fmt.Errorf("v1ファイルの書き込みに失敗しました。(Failed to write v1 file.) [%s]", err)
			}
		}
	}
	return nil
}

func DownloadBackgroundGenerator(destPath string) error {
	const downloadURL = "https://github.com/TootieJin/pjsekai-background-gen-rust/releases/download/v0.1.0/pjsekai-background-gen.exe"

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("ダウンロードサーバーに接続できませんでした。(Could not connect to download server.) [%s]", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ファイルのダウンロードに失敗しました。(Failed to download file.) [%d]", resp.StatusCode)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました。(Failed to create file.) [%s]", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました。(Failed to write file.) [%s]", err)
	}

	return nil
}
