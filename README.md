[**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール-日本語)

[![Releases](https://img.shields.io/github/downloads/TootieJin/pjsekai-overlay-APPEND/total.svg)](https://gitHub.com/TootieJin/pjsekai-overlay-APPEND/releases/)
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool (English)
Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com), an open-sourced tool to make Project Sekai Fanmade (custom chart) videos.

- **16:9**

https://github.com/user-attachments/assets/dda7225a-a7f3-41d4-bbf4-9cec9b03b840

- **4:3 (Tournament Mode ON)**

https://github.com/user-attachments/assets/ab4ee52c-2ffa-4941-b916-87e1f3559d72

- **v1 Skin (1e+30 power)**

https://github.com/user-attachments/assets/3efab743-246a-4da7-8d80-a02b2f09f5b3

- **Video Example**

*(Click the image to watch it)*\
[![【Project Sekai x Honkai: Star Rail】Nameless Faces - HOYO-MiX feat. Lilas Ikuta (Fanmade)](https://img.youtube.com/vi/uXx1OZDQZOI/maxresdefault.jpg)](https://youtu.be/uXx1OZDQZOI)
[![【Project Sekai Fanmade? (v3→v1)】Hello, SEKAI - DECO*27【ETERNAL Lv32】](https://img.youtube.com/vi/BHVNuwxA1ek/maxresdefault.jpg)](https://youtu.be/BHVNuwxA1ek)

> [!CAUTION]
> **For English users:** This tool is primary only for people with technical know-how and basic knowledge of AviUtl.\
> Only use this tool if you can figure it out yourself. **DO NOT open issues, DM me, or request help in Sonolus / Chart Cyanvas Discord servers about this**.

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - [Extra assets](./assets/extra%20assets) (thank you [ReiyuN](https://discordid.netlify.app/?id=383636820409188374), [Gaven](https://github.com/gaven1880) and [YumYummity](https://github.com/YumYummity) for the contribution!)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos
  - iPad (4:3) video support
  - Ability to use the English AviUtl
  - v1 UI skin (Full support)
  - Automatically changes chart difficulty to generate in AviUtl based on chart tag _(or title)_
  - Increased score limit to infinity (?)
  [![image](https://github.com/user-attachments/assets/baceaf22-fdcb-4b48-8fb7-54b08e6d3086)]()
  [![pjsekai-overlay-APPEND_minint](https://github.com/user-attachments/assets/80eb8fc1-6602-4c26-ac47-4e8e07fb99c2)]()
  [![pjsekai-overlay-APPEND_maxint](https://github.com/user-attachments/assets/45a49c19-7883-402c-b016-58f02f72f0b6)]()


## Terms of Use

**(REQUIRED)** In the description of your video, please copy the text here:

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
   https://github.com/TootieJin/pjsekai-overlay-MULTI (Multilive ver.)
- Original by 名無し｡ (https://sevenc7c.com) 
   https://github.com/sevenc-nanashi/pjsekai-overlay
```

**JP**
```
プロセカ風動画作成補助ツール：
- TootieJin (https://tootiejin.com) フォーク版
   https://github.com/TootieJin/pjsekai-overlay-APPEND
   https://github.com/TootieJin/pjsekai-overlay-MULTI (マルチライブver.)
- 名無し｡ (https://sevenc7c.com) オリジナル版
   https://github.com/sevenc-nanashi/pjsekai-overlay
```

> [!NOTE]
> **(optional)** You can remove watermark by check/unchecking `Watermark` in the `Root@pjsekai-overlay-en` element.\
> ![image](https://github.com/user-attachments/assets/48a636da-8ec0-443b-9cf4-b73fd93c47df)

## Requirements

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest)
  - (Recommended: [patch.aul](https://scrapbox.io/ePi5131/patch.aul))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- Basic knowledge of AviUtl

*- Refer to this [English guide](https://github.com/Khronophobia/pjsekai-overlay-english/wiki/Usage-Guide) on how to use AviUtl EN.*
> [!IMPORTANT]
> **REMEMBER TO GO TO `File > ENVIRONMENT SETTINGS > SYSTEM SETTINGS` AND SET THE `Max image size` TO 4000x3000 (or bigger)!!!!!!!!**

## Video Guide

1. [Make your chart first.](https://cc.sevenc7c.com)
2. Go to [Sonolus](https://sonolus.com/) to find your chart.
3. Screen record the video with **BLACK background** and「Hide UI」turned on
4. Transfer the video file to your computer.
   - Download the [ffmpeg](https://www.ffmpeg.org/) encoder if you haven't.
5. Once done, refer to the usage guide below.

## Usage Guide

1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
2. Unzip the file
3. Import ALL files & folders in the [`depenencies/aviutl script`](./depenencies/aviutl%20script) folder into here:
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *If a folder is missing, make a new folder with said name.*
4. Open AviUtl
   - **Note: You must open AviUtl before opening pjsekai-overlay-APPEND to install objects.**
5. Open `pjsekai-overlay-APPEND.exe`
6. Input the chart ID including the prefix.
   - `chcy-`: Chart Cyanvas (cc.sevenc7c.com)
   - `ptlv-`: Potato Leaves (ptlv.sevenc7c.com)
   - `utsk-`: Untitled Sekai (us.pim4n-net.com)
7. Import specified exo file by navigating to your `pjsekai-overlay/dist/[chart ID]` directory:
   - **For phone users:** main_en_16-9_1920x1080.exo
   - **For iPad users:** main_en_4-3_1440x1080.exo

---------------------------------------------------------------------------------------

# pjsekai-overlay-APPEND / フォークプロセカ風動画作成補助ツール (日本語)

[TootieJin](https://tootiejin.com)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)用フォーク。
pjsekai-overlay(-APPEND) は、プロセカの創作譜面をプロセカ風の動画にするためのオープンソースのツールです。

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - [追加アセット](./assets/extra%20assets/) ([ReiyuN](https://discordid.netlify.app/?id=383636820409188374)さん、[Gaven](https://github.com/gaven1880)さんと[YumYummity](https://github.com/YumYummity)さん、ご寄稿ありがとうございました。)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成
  - iPad（4:3）動画対応
  - 英語版AviUtlの使用機能
  - v1 UIスキン（フル対応）
  - 譜面のタグ _（またはタイトル）_ に基づいて、AviUtlで生成される譜面の難易度を自動的に変更する
  - スコアの上限を無限大（？）
  [![image](https://github.com/user-attachments/assets/baceaf22-fdcb-4b48-8fb7-54b08e6d3086)]()
  [![pjsekai-overlay-APPEND_minint](https://github.com/user-attachments/assets/80eb8fc1-6602-4c26-ac47-4e8e07fb99c2)]()
  [![pjsekai-overlay-APPEND_maxint](https://github.com/user-attachments/assets/45a49c19-7883-402c-b016-58f02f72f0b6)]()


## 利用規約

**(必須)** 動画の説明文に、こちらのテキストをコピーしてください：

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
   https://github.com/TootieJin/pjsekai-overlay-MULTI (Multilive ver.)
- Original by 名無し｡ (https://sevenc7c.com) 
   https://github.com/sevenc-nanashi/pjsekai-overlay
```

**JP**
```
プロセカ風動画作成補助ツール：
- TootieJin (https://tootiejin.com) フォーク版
   https://github.com/TootieJin/pjsekai-overlay-APPEND
   https://github.com/TootieJin/pjsekai-overlay-MULTI (マルチライブver.)
- 名無し｡ (https://sevenc7c.com) オリジナル版
   https://github.com/sevenc-nanashi/pjsekai-overlay
```

> [!NOTE]
> **(任意)** `設定@pjsekai-overlay`要素でチェック/チェックを外すことで、`透かし`を消すことができます。\
> ![image](https://github.com/user-attachments/assets/05cc6e7d-0e62-4729-a4df-c8d8634b0a10)

## 必須事項

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) （[導入方法](https://aviutl.info/dl-innsuto-ru/)）
  - (強く推奨：[patch.aul](https://scrapbox.io/ePi5131/patch.aul))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- AviUtlの基本的な知識

## 動画の作り方

1. [譜面を作る](https://wiki.purplepalette.net/create-charts)
2. [Sonolus](https://sonolus.com/)で譜面を撮影する
3. **背景を黒**にし、「Hide UI」をONにして、動画をスクリーン録画します。
4. 撮影したプレイ動画のファイルをパソコンに転送する
   - Google Drive など
5. [ffmpeg](https://www.ffmpeg.org/)で再エンコードする
   - AviUtl で読み込むため
6. 下の利用方法に従って UI を後付けする

## 利用方法

1. 右の Releases から最新のバージョンの zip をダウンロードする
2. zip を解凍する
3. [`depenencies/aviutl script`](./depenencies/aviutl%20script)フォルダ内のすべてのファイルとフォルダをここにインポートします：
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *フォルダがない場合は、その名前で新しいフォルダを作ってください。*
4. AviUtl を起動する
   - **pjsekai-overlay が起動する前に AviUtl を起動するとオブジェクトのインストールが行われます。**
5. `pjsekai-overlay-APPEND.exe` を起動する
6. 譜面IDを接頭辞込みで入力して下さい
   - `chcy-`: Chart Cyanvas (cc.sevenc7c.com)
   - `ptlv-`: Potato Leaves (ptlv.sevenc7c.com)
   - `utsk-`: Untitled Sekai (us.pim4n-net.com)
7. `pjsekai-overlay/dist/[譜面ID]`ディレクトリに移動して、指定したexoファイルをインポートします：
   - **スマホ向け:** main_jp_16-9_1920x1080.exo
   - **iPad向け:** main_jp_4-3_1440x1080.exo
