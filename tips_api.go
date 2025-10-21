package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type Tip struct {
	JapaneseText string `json:"japanese_text"`
	EnglishText  string `json:"english_text"`
}

type TipsClient struct {
	BaseURL string
	Client  *http.Client
}

func NewTipsClient(baseURL string) *TipsClient {
	return &TipsClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (tc *TipsClient) GetRandomTip() (*Tip, error) {
	resp, err := tc.Client.Get(tc.BaseURL + "/api/tip/random")
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
	}

	var tip Tip
	if err := json.Unmarshal(body, &tip); err != nil {
		return nil, fmt.Errorf("JSONパースエラー: %v", err)
	}

	return &tip, nil
}

func TipsAPI() {
	client := NewTipsClient("https://overlay-tips.pim4n.f5.si/")

	tip, err := client.GetRandomTip()
	if err != nil {
		fmt.Printf(color.CyanString("◆ APIからのTip取得に失敗しました: %v\n", err))
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tips := []string{
			"水板を継続して運用する場合、名無しは月額15,000円以上を支払う必要がありました。",
			"Had Chart Cyanvas server continue to run, Nanashi would have to pay $100+/month for it.",

			"overlayのアップデート後は、objファイルを再インストールすることを忘れないでください！",
			"Remember to reinstall object file when you update overlay!",

			"v1のUIが懐かしい…",
			"I miss v1 UI...",

			"1 .",
			"█   .     3.....",
		}

		a := r.Intn(len(tips) - 1)
		if a%2 != 0 {
			a--
		}
		fmt.Printf(color.CyanString("◆ Tip: %s\n◆ Tip: %s\n\n"), tips[a], tips[a+1])
		return
	}

	fmt.Printf(color.CyanString("◆ Tip: %s\n◆ Tip: %s\n\n"), tip.JapaneseText, tip.EnglishText)
}
