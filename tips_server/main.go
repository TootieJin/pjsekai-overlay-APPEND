package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Tip struct {
	JapaneseText string `json:"japanese_text"`
	EnglishText  string `json:"english_text"`
}

var allTips []Tip

func init() {
	allTips = []Tip{
		{
			JapaneseText: "水板を継続して運用する場合、名無しは月額15,000円以上を支払う必要がありました。",
			EnglishText:  "Had Chart Cyanvas server continue to run, Nanashi would have to pay $100+/month for it.",
		},
		{
			JapaneseText: "overlayのアップデート後は、objファイルを再インストールすることを忘れないでください！",
			EnglishText:  "Remember to reinstall object file when you update overlay!",
		},
		{
			JapaneseText: "「extra assets」フォルダーの中に、一体どんな秘密が隠されているのか...",
			EnglishText:  "Who knows what secret things lie in the \"extra assets\" folder...",
		},
		{
			JapaneseText: "「data.ped」ファイルでスコア、コンボ、その他の項目を編集できます。",
			EnglishText:  "You can edit score, combo & other things in the \"data.ped\" file.",
		},
		{
			JapaneseText: "!!",
			EnglishText:  "!!",
		},
		{
			JapaneseText: "??",
			EnglishText:  "??",
		},
		{
			JapaneseText: "Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip:",
			EnglishText:  "Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip:",
		},
		{
			JapaneseText: "「APPEND」とは、ゲームプレイのために別の指を追加(APPEND)することを意味します。",
			EnglishText:  "[APPEND] means you APPEND another finger to play the game.",
		},
		{
			JapaneseText: "AviUtlは動画編集ツールです。UIで様々なカスタマイズが可能です。",
			EnglishText:  "AviUtl is an video editing tool. You can do crazy things with the UI.",
		},
		{
			JapaneseText: "このTipはDeepLを使用して翻訳されています。",
			EnglishText:  "The tip above uses DeepL to translate this tip.",
		},
		{
			JapaneseText: "画像が切り取られていますか？「ファイル > 設定 > システム」で「最大画像サイズ」を設定する必要があります。",
			EnglishText:  "Image cropped? You need to adjust \"Max resolution\" in \"File > SETTINGS > SYSTEM\".",
		},
		{
			JapaneseText: "どれくらいのTipsが実際に役立つと思いますか？",
			EnglishText:  "Guess how many Tips are literally useful?",
		},
		{
			JapaneseText: "この2年半、Chart Cyanvasをご利用いただき、誠にありがとうございました。",
			EnglishText:  "Thank you all for being with Chart Cyanvas for the past 2.5 years.",
		},
		{
			JapaneseText: "AviUtlに元の動画ファイルをインポートすると、動画が同期しなくなります。FFmpegを使用して、そこでエンコードしてください。",
			EnglishText:  "Importing raw video file to AviUtl makes the video out of sync. Use FFmpeg and encode it there.",
		},
		{
			JapaneseText: "pjsekai-overlay-APPENDを使用すると、TikTokの子供をだまして本物だと信じ込ませることができます。（推奨されません）",
			EnglishText:  "pjsekai-overlay-APPEND can be used to fool TikTok kids into thinking it's real. (Not recommended)",
		},
		{
			JapaneseText: "プロセカで使用されているフォントは、「dependencies」フォルダー内に格納されています。",
			EnglishText:  "The fonts used in sekai can be found in the \"dependencies\" folder...",
		},
		{
			JapaneseText: "「data.ped」ファイルの各行は、以下の形式に従っています：s|[時間枠（秒）]:[合計スコア]:[追加スコア]:[スコアバーの位置]:[順位]:[コンボ]",
			EnglishText:  "Each line in the \"data.ped\" file follows this format: s|[timeframe(sec)]:[totalscore]:[addedscore]:[scorebar position]:[rank]:[combo]",
		},
		{
			JapaneseText: "APコンボ判定は、AviUtlにおいて互換性があります。",
			EnglishText:  "AP Combo & Judgement can be interchangeable in AviUtl.",
		},
		{
			JapaneseText: "設定@pjsekai-overlayでのオフセットはあなたの味方です。",
			EnglishText:  "OFFSET in Root@pjsekai-overlay-en is your friend.",
		},
		{
			JapaneseText: "総合力を250000に設定する代わりに、なぜもっと高くしないのか？無限大まで、つまり。",
			EnglishText:  "Instead of setting team power at 250000, why not go higher? To infinity, that is.",
		},
		{
			JapaneseText: "総合力をマイナス数値に設定できます。試してみてください。",
			EnglishText:  "You can set team power to a negative number. Try it.",
		},
		{
			JapaneseText: "古いUIが恋しいですか？お任せください。",
			EnglishText:  "Miss the old UI? I'm here for you.",
		},
		{
			JapaneseText: "フリックの遊び方？ ↑←→↗↖↗↑→←",
			EnglishText:  "How to play Flick notes? ↑←→↗↖↗↑→←",
		},
		{
			JapaneseText: "公式エンジンでは、最後の4桁のコンボ番号のみが表示されます。（例：12345 → █2345）",
			EnglishText:  "In the official engine, only the last 4 combo digits are rendered. (e.g. 12345 → █2345)",
		},
		{
			JapaneseText: "注意！前回のTipを覚えていますか？",
			EnglishText:  "Attention! Do you remember last tip?",
		},
		{
			JapaneseText: "すべて「愛おしい」と思う季節にさよなら (Say goodbye when everything is \"lovely\")",
			EnglishText:  "いまだけは (Just for now)",
		},
		{
			JapaneseText: "ああ、このTipは話題がそれました。すみません。",
			EnglishText:  "Ah, this tip went off-topic. Sorry.",
		},
		{
			JapaneseText: "Tipが見つかりません。",
			EnglishText:  "Tip not found.",
		},
		{
			JapaneseText: "設定@pjsekai-overlay要素でチェックを外すことで、「透かし」を消すことができます。",
			EnglishText:  "You can remove watermark by unchecking \"Watermark\" in the Root@pjsekai-overlay-en element.",
		},
		{
			JapaneseText: "リポジトリに追加したい別の水板インスタンスがありますか？PRを作成してください。",
			EnglishText:  "Have another Chart Cyanvas instance you want to add in the repo? Make a pull request.",
		},
		{
			JapaneseText: "[非公開]",
			EnglishText:  "[REDACTED]",
		},
		{
			JapaneseText: "プロセカについて、関係のない場面で言及しないよう、よろしいでしょうか？",
			EnglishText:  "Would you mind not mentioning Project Sekai on irrelevant occasions?",
		},
		{
			JapaneseText: "ここで何を書けばいいか、少し考えてみます。",
			EnglishText:  "Let me think what I should write here.",
		},
		{
			JapaneseText: "AUTOLIVEはどこですか？ 画面の右下にあります。",
			EnglishText:  "Where's the auto? It's at the bottom right corner.",
		},
		{
			JapaneseText: "ザ",
			EnglishText:  "The",
		},
		{
			JapaneseText: "このTipが表示されている場合は、無視して構いません。",
			EnglishText:  "Japanese characters look like gibberish? Go to your language settings, Administrative language settings and change the Language for non-Unicode programs to Japanese.",
		},
		{
			JapaneseText: "最も難しい創作譜面は、ブラックホールのように魅力的です。",
			EnglishText:  "The hardest custom charts are as attractive as black holes.",
		},
		{
			JapaneseText: "一部の譜面作成者は、本日「ブルーアーカイブ」をプレイしています。",
			EnglishText:  "Some charters are playing this cunny game \"Blue Archive\" today.",
		},
		{
			JapaneseText: "一部の譜面作成者は、本日「ウマ娘」をプレイしています。",
			EnglishText:  "Some charters are playing the horse game \"Umamusume\" today.",
		},
		{
			JapaneseText: "38面ダイスを使って、公式譜面の難易度を決定します。",
			EnglishText:  "A 38-sided dice are used to decide the difficulty of each official chart.",
		},
		{
			JapaneseText: "セガ（英語）は、「Anime Expo 2025 × プロセカ(EN)」キャンペーンを実施し、特定の曲を100万回プレイすることで...300クリスタルを獲得できるキャンペーンを実施しました。",
			EnglishText:  "SEGA (English) hosted an \"Anime Expo 2025 x Colorful Stage\" campaign requiring everyone to play a specific song 1 MILLION times to get... 300 crystals.",
		},
		{
			JapaneseText: "\n                            —{›\n                           —íí{\n                    —{{›   —íí{    {{{\n                   —íííí›  —íí{   {íííí\n                   íííí{   —íí{   —íííí›\n                  —íííí—   —íí{    íííí{\n      ››››››››››››ííííí››››—íí{››››—íííí—›››››››››››\n    ›íííííííííííííííííííííííííííííííííííííííííííííííí{\n    ííí———íí———íííí———————————————————————————ííí——ííí›\n    ííí› ›íí››í—                              {í{  {íí›\n    ííí› ›íí›íííííííí—             ííííííííííííí{  {íí›\n    ííí› ›íí› › › ›››—{{{{{{—›      › › › › › íí{  {íí›\n    ííí› ›íí›            ›{{{{{{{—›           íí{  {íí›\n    ííí› ›íí›             —{{{{{{{›     ›—{›  {í{  {íí›\n    ííí› ›íí›         ››—{{{{{——›      ›—››í— íí{  {íí›\n    ííí› ›íí—íííííííí{——››           —ííííííí{íí{  {íí›\n    ííí› ›íí››—›—›—››       ››——{{{{{———›—›—››{í{  {íí›\n    ííí› ›íí›           ›{{{{{{{—›            íí{  {íí›\n    ííí› ›íí›          ›{{{{{{{{              íí{  {íí›\n    ííí› ›íí›             —{{{{{{{—           {í{  {íí›\n    ííí› ›íí—íííííííííí{        ››—{{íííííííí{íí{  {íí›\n    ííí› ›íí›———————————             ›———————›íí{  {íí›\n    ííí› ›íí›››››››››                         {í{  {íí›\n {ííííííííííííííííííííííííííííííííííííííííííííííííííííííí›\n›íííííííííííííííííííííííííííííííííííííííííííííííííííííííí{\n        {íííí›             —íí{              {íííí\n       ›íííí—              —íí{              ›íííí{\n       {íííí›              —íí{               {íííí›\n      ›íííí{               —íí{               ›íííí—\n      {íííí                —íí{                {íííí›\n     —íííí—                —íí{                ›íííí{\n     {ííííí{{{{{{{{{{{{{{{{íííí{{{{{{{{{{{{{{{{{ííííí\n    ›íííííííííííííííííííííííííííííííííííííííííííííííí{\n    ííííí                                        {íííí›\n   —íííí—                                        ›íííí{\n   {íííí                                          {íííí›\n  —íííí—                                           íííí{\n  {íííí                                            {íííí›\n ›íííí—                                            ›íííí—\n  ›íí—                                              ›{í—",
			EnglishText:  "Chart Cyanvas",
		},
		{
			JapaneseText: "このTipは、█回に1回表示されます ￣︶￣",
			EnglishText:  "This tip will be displayed once every █ times ￣︶￣",
		},
		{
			JapaneseText: "一部のTipはPhigrosから借用しています。なぜなら、私たちのクリエイターは創造力が限られているからです。",
			EnglishText:  "Some tips are stolen from Phigros, because our creator has limited creativity.",
		},
		{
			JapaneseText: "プログラムを使用する際、NaNエラーが発生しています。",
			EnglishText:  "I have found NaN errors when you use the program.",
		},
		{
			JapaneseText: "譜面作成者として人気を得るには？TikTokで流行りの曲を譜面化すればOK！",
			EnglishText:  "How to be popular as a charter? Just chart a song that's trending on TikTok.",
		},
		{
			JapaneseText: "README を読む前に、使えないことに怒らないでください。本当です。読んでください。",
			EnglishText:  "Read README before being mad that you can't use it. No, really. Read it.",
		},
		{
			JapaneseText: "CC分岐サーバー（chart-cyanvas.com）は2025年9月13日に作成された。今日に至るまで、誰がサイトをホストしているのか誰も知らなかった。",
			EnglishText:  "Chart Cyanvas Offshoot Server (chart-cyanvas.com) was made in September 13th, 2025. To this day, nobody knew who hosted the site.",
		},
		{
			JapaneseText: "v1のUIが懐かしい…",
			EnglishText:  "I miss v1 UI...",
		},
		{
			JapaneseText: "1 .",
			EnglishText:  "█   .     3.....",
		},
	}
}

func main() {
	// ランダムなtipを返すAPIエンドポイント
	http.HandleFunc("/api/tip/random", func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		tipIndex := rand.Intn(len(allTips))
		tip := allTips[tipIndex]

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // CORSヘッダー
		json.NewEncoder(w).Encode(tip)
	})

	// 全てのtipsを返すAPIエンドポイント
	http.HandleFunc("/api/tips", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // CORSヘッダー
		json.NewEncoder(w).Encode(allTips)
	})

	// ヘルスチェック用のエンドポイント
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is running")
	})

	// サーバー起動
	port := "8050"
	fmt.Printf("Tips APIサーバーが起動しました。ポート: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
