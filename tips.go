package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// これを見て何か追加したいTipsがあれば、PRを送ってください
// if you see this and you have some tips you would like to add, make a PR pls

func Tips() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tips := []string{
		"水板を継続して運用する場合、名無しは月額15,000円以上を支払う必要がありました。",
		"Had Chart Cyanvas server continue to run, Nanashi would have to pay $100+/month for it.",

		"overlayのアップデート後は、objファイルを再インストールすることを忘れないでください！",
		"Remember to reinstall object file when you update overlay!",

		"「extra assets」フォルダーの中に、一体どんな秘密が隠されているのか...",
		"Who knows what secret things lie in the \"extra assets\" folder...",

		"「data.ped」ファイルでスコア、コンボ、その他の項目を編集できます。",
		"You can edit score, combo & other things in the \"data.ped\" file.",

		"!!",
		"!!",

		"??",
		"??",

		"Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip:",
		"Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip: Tip:",

		"「APPEND」とは、ゲームプレイのために別の指を追加(APPEND)することを意味します。",
		"[APPEND] means you APPEND another finger to play the game.",

		"AviUtlは動画編集ツールです。UIで様々なカスタマイズが可能です。",
		"AviUtl is an video editing tool. You can do crazy things with the UI.",

		"このTipはDeepLを使用して翻訳されています。",
		"The tip above uses DeepL to translate this tip.",

		"画像が切り取られていますか？「ファイル > 設定 > システム」で「最大画像サイズ」を設定する必要があります。",
		"Image cropped? You need to adjust \"Max resolution\" in \"File > SETTINGS > SYSTEM\".",

		"どれくらいのTipsが実際に役立つと思いますか？",
		"Guess how many Tips are literally useful?",

		"この2年半、Chart Cyanvasをご利用いただき、誠にありがとうございました。",
		"Thank you all for being with Chart Cyanvas for the past 2.5 years.",

		"AviUtlに元の動画ファイルをインポートすると、動画が同期しなくなります。FFmpegを使用して、そこでエンコードしてください。",
		"Importing raw video file to AviUtl makes the video out of sync. Use FFmpeg and encode it there.",

		"pjsekai-overlay-APPENDを使用すると、TikTokの子供をだまして本物だと信じ込ませることができます。（推奨されません）",
		"pjsekai-overlay-APPEND can be used to fool TikTok kids into thinking it's real. (Not recommended)",

		"プロセカで使用されているフォントは、「dependencies」フォルダー内に格納されています。",
		"The fonts used in sekai can be found in the \"dependencies\" folder...",

		"「data.ped」ファイルの各行は、以下の形式に従っています：s|[時間枠（秒）]:[合計スコア]:[追加スコア]:[スコアバーの位置]:[順位]:[コンボ]",
		"Each line in the \"data.ped\" file follows this format: s|[timeframe(sec)]:[totalscore]:[addedscore]:[scorebar position]:[rank]:[combo]",

		"APコンボ判定は、AviUtlにおいて互換性があります。",
		"AP Combo & Judgement can be interchangeable in AviUtl.",

		"設定@pjsekai-overlayでのオフセットはあなたの味方です。",
		"OFFSET in Root@pjsekai-overlay-en is your friend.",

		"総合力を250000に設定する代わりに、なぜもっと高くしないのか？無限大まで、つまり。",
		"Instead of setting team power at 250000, why not go higher? To infinity, that is.",

		"総合力をマイナス数値に設定できます。試してみてください。",
		"You can set team power to a negative number. Try it.",

		"古いUIが恋しいですか？お任せください。",
		"Miss the old UI? I'm here for you.",

		"フリックの遊び方？ ↑←→↗↖↗↑→←",
		"How to play Flick notes? ↑←→↗↖↗↑→←",

		"公式エンジンでは、最後の4桁のコンボ番号のみが表示されます。（例：12345 → █2345）",
		"In the official engine, only the last 4 combo digits are rendered. (e.g. 12345 → █2345)",

		"注意！前回のTipを覚えていますか？",
		"Attention! Do you remember last tip?",

		"すべて「愛おしい」と思う季節にさよなら (Say goodbye when everything is \"lovely\")",
		"いまだけは (Just for now)", // Ref:rain

		"ああ、このTipは話題がそれました。すみません。",
		"Ah, this tip went off-topic. Sorry.", // this is the reason why Sonolus Discord server crashed w

		"Tipが見つかりません。",
		"Tip not found.",

		"設定@pjsekai-overlay要素でチェックを外すことで、「透かし」を消すことができます。",
		"You can remove watermark by unchecking \"Watermark\" in the Root@pjsekai-overlay-en element.",

		"リポジトリに追加したい別の水板インスタンスがありますか？PRを作成してください。",
		"Have another Chart Cyanvas instance you want to add in the repo? Make a pull request.",

		"[非公開]",
		"[REDACTED]",

		"プロセカについて、関係のない場面で言及しないよう、よろしいでしょうか？",
		"Would you mind not mentioning Project Sekai on irrelevant occasions?",

		"ここで何を書けばいいか、少し考えてみます。",
		"Let me think what I should write here.",

		"AUTOLIVEはどこですか？ 画面の右下にあります。",
		"Where's the auto? It's at the bottom right corner.", // TikTok comment trend

		"ザ",
		"The",

		"このTipが表示されている場合は、無視して構いません。",
		"Japanese characters look like gibberish? Go to your language settings, Administrative language settings and change the Language for non-Unicode programs to Japanese.",

		"最も難しい創作譜面は、ブラックホールのように魅力的です。",
		"The hardest custom charts are as attractive as black holes.",

		"一部の譜面作成者は、本日「ブルーアーカイブ」をプレイしています。",
		"Some charters are playing this cunny game \"Blue Archive\" today.",

		"一部の譜面作成者は、本日「ウマ娘」をプレイしています。",
		"Some charters are playing the horse game \"Umamusume\" today.",

		"一部の譜面作成者は、本日「ステラソラ」をプレイしています。",
		"Some charters are playing \"Stella Sora\" today.",

		"一部の譜面作成者は、本日「ユメステ」をプレイしています。",
		"Some charters are playing the harder game \"World Dai Star\" today.",

		"38面ダイスを使って、公式譜面の難易度を決定します。",
		"A 38-sided dice are used to decide the difficulty of each official chart.",

		"ああ…38面ダイスは38を出した…",
		"Oh...  the 38-sided dice landed a 38...",

		"セガ（英語）は、「Anime Expo 2025 × プロセカ(EN)」キャンペーンを実施し、特定の曲を100万回プレイすることで...300クリスタルを獲得できるキャンペーンを実施しました。",
		"SEGA (English) hosted an \"Anime Expo 2025 x Colorful Stage\" campaign requiring everyone to play a specific song 1 MILLION times to get... 300 crystals.",

		"\n                            —{›\n                           —íí{\n                    —{{›   —íí{    {{{\n                   —íííí›  —íí{   {íííí\n                   íííí{   —íí{   —íííí›\n                  —íííí—   —íí{    íííí{\n      ››››››››››››ííííí››››—íí{››››—íííí—›››››››››››\n    ›íííííííííííííííííííííííííííííííííííííííííííííííí{\n    ííí———íí———íííí———————————————————————————ííí——ííí›\n    ííí› ›íí››í—                              {í{  {íí›\n    ííí› ›íí›íííííííí—             ííííííííííííí{  {íí›\n    ííí› ›íí› › › ›››—{{{{{{—›      › › › › › íí{  {íí›\n    ííí› ›íí›            ›{{{{{{{—›           íí{  {íí›\n    ííí› ›íí›             —{{{{{{{›     ›—{›  {í{  {íí›\n    ííí› ›íí›         ››—{{{{{——›      ›—››í— íí{  {íí›\n    ííí› ›íí—íííííííí{——››           —ííííííí{íí{  {íí›\n    ííí› ›íí››—›—›—››       ››——{{{{{———›—›—››{í{  {íí›\n    ííí› ›íí›           ›{{{{{{{—›            íí{  {íí›\n    ííí› ›íí›          ›{{{{{{{{              íí{  {íí›\n    ííí› ›íí›             —{{{{{{{—           {í{  {íí›\n    ííí› ›íí—íííííííííí{        ››—{{íííííííí{íí{  {íí›\n    ííí› ›íí›———————————             ›———————›íí{  {íí›\n    ííí› ›íí›››››››››                         {í{  {íí›\n {ííííííííííííííííííííííííííííííííííííííííííííííííííííííí›\n›íííííííííííííííííííííííííííííííííííííííííííííííííííííííí{\n        {íííí›             —íí{              {íííí\n       ›íííí—              —íí{              ›íííí{\n       {íííí›              —íí{               {íííí›\n      ›íííí{               —íí{               ›íííí—\n      {íííí                —íí{                {íííí›\n     —íííí—                —íí{                ›íííí{\n     {ííííí{{{{{{{{{{{{{{{{íííí{{{{{{{{{{{{{{{{{ííííí\n    ›íííííííííííííííííííííííííííííííííííííííííííííííí{\n    ííííí                                        {íííí›\n   —íííí—                                        ›íííí{\n   {íííí                                          {íííí›\n  —íííí—                                           íííí{\n  {íííí                                            {íííí›\n ›íííí—                                            ›íííí—\n  ›íí—                                              ›{í—",
		"Chart Cyanvas", // ASCII art

		"このTipは、█回に1回表示されます ￣︶￣",
		"This tip will be displayed once every █ times ￣︶￣",

		"一部のTipはPhigrosから借用しています。なぜなら、私たちのクリエイターは創造力が限られているからです。",
		"Some tips are stolen from Phigros, because our creator has limited creativity.",

		"プログラムを使用する際、NaNエラーが発生しています。",
		"I have found NaN errors when you use the program.",

		"譜面作成者として人気を得るには？TikTokで流行りの曲を譜面化すればOK！",
		"How to be popular as a charter? Just chart a song that's trending on TikTok.",

		"README を読む前に、使えないことに怒らないでください。本当です。読んでください。",
		"Read README before being mad that you can't use it. No, really. Read it.",

		"CC分岐サーバー（chart-cyanvas.com）は2025年9月13日に作成された。今日に至るまで、誰がサイトをホストしているのか誰も知らなかった。",
		"Chart Cyanvas Offshoot Server (chart-cyanvas.com) was made in September 13th, 2025. To this day, nobody knew who hosted the site.",

		"v1のUIが懐かしい…",
		"I miss v1 UI...",

		"1 .",
		"█   .     3.....",

		"sudo device auto-play",
		"sudo device auto-play",

		"← To Be Continued...",
		"← To Be Continued...",

		"怪獣になりたい",
		"I want to be a monster",

		"[非表示のTip]",
		"[Hidden Tip]",

		"▼・ᴗ・▼",
		"▼・ᴗ・▼",

		"ミクさと一歌さはいつも君と一緒にいる…プロセカをアンインストールしない限り！",
		"Miku and Ichika will always be with you... Unless you uninstall the game!",

		"pjsekai-overlay-APPENDが1周年を迎えます（2024年11月24日）！本ツールをご利用いただきありがとうございます。",
		"pjsekai-overlay-APPEND is turning a year old (Nov 24th, 2024)! Thank you for using the tool.",

		"Don't Fight The Music",
		"And Revive The Melody",

		"これって本当に読まれてるの？",
		"Do people actually read these?",

		"今日に至るまで、どちらのAviUtlバージョンにもmp4出力機能が組み込まれていない...",
		"To this day, both AviUtl versions doesn't even have built-in mp4 exporter...",

		"* 行動ごとにノーツが評価されます。獲得した経験値ごとにノーツが評価されます。",
		"* Notes will be judged for every action. Notes will be judged for every EXP you've earned.",

		"この番組の制作資金は、視聴者の皆様のご支援により実現しました",
		"Funding for this program was made possible by viewers like you",

		"知ってる？読み込み中に何か読みたいから、Tipがあるんだよ…",
		"Do you know? Tips exist cuz you want something to read when loading...",

		"このツールが好きならリポジトリをスターしてください！（あるいは好きじゃなくても！）",
		"Star the repo if you love this tool! (or don't love this tool!)",

		"ここをクリックして素早く生成",
		"Click here to quick-generate",

		"AviUtl2に英語モジュールがないということは、別途の英語エイリアス(.object)ファイルを使用する必要がないことを意味します。",
		"AviUtl2 not having English mod means you don't have to use the separate English alias(.object) file.",

		"ドラッグ＆ドロップ時にアニメーションが遅くなることに気づきましたか？「プロジェクトの新規」でフレームレートを60fpsに設定してください。",
		"Noticed the animation slowing down when drag and dropped? Set frame rate in \"New Project\" to 60fps.",

		"外の月を見て！今すぐ！",
		"Look at the moon outside! RIGHT NOW!",

		"これまでにない譜面を作成しよう",
		"Chart like you never did before",

		"PLAY    YOU     DID",
		"    LIKE   NEVER   BEFORE",

		"最適な体験のためにはAviUtl ExEdit2をご利用ください。",
		"Use AviUtl ExEdit2 for best experience.",

		"2025年11月8日現在、公式ゲーム内の譜面総数は3107枚である。二項分布の公式を用いると、38面ダイスでLv38が出る確率は1.1×10^34分の1である。",
		"As of Nov 8th 2025, there are a grand total of 3107 charts in the official game. Using the binomial formula, there's a 1 in 1.1×10^34 chance land a Lv38 in a 38-sided die.",
	}

	// これを見て何か追加したいTipsがあれば、PRを送ってください
	// if you see this and you have some tips you would like to add, make a PR pls

	a := r.Intn(len(tips) - 1)
	if a%2 != 0 {
		a--
	}
	fmt.Printf(color.CyanString("◆ Tip: %s\n◆ Tip: %s\n\n"), tips[a], tips[a+1])
}
