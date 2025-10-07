[**English Section**](#pjsekai-overlay-append--forked-pjsekai-style-video-creation-tool-english) | [**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール-日本語)

[![Releases](https://img.shields.io/github/downloads/TootieJin/pjsekai-overlay-APPEND/total.svg)](https://gitHub.com/TootieJin/pjsekai-overlay-APPEND/releases/)
[![pjsekai-overlay-APPEND thumbnail](https://github.com/user-attachments/assets/dcc037f7-1c2b-4b83-b17b-b3c5155d670c)]()
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool (English)
Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com), an open-sourced tool to make Project Sekai Fanmade (custom chart) videos.

> [!CAUTION]
> **For English users:** This tool is primary only for people with technical know-how and basic knowledge of AviUtl.\
> Only use this tool if you can figure it out yourself. **DO NOT open issues, DM me, or request help about this**.

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

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - [Extra assets](./assets/extra%20assets) (thank you [ReiyuN](https://github.com/Reiyunkun), [Gaven](https://github.com/gaven1880) and [YumYummity](https://github.com/YumYummity) for the contribution!)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos
  - iPad (4:3) video support
  - Ability to use the English AviUtl
  - v1 UI skin (Full support)
  - Automatically changes chart difficulty to generate in AviUtl based on chart tag (or title)
  - Additional support for more servers
  - Various UI customization
    - Animated Scoring
    - Toggle on/off animation in different elements
    - Interchangable AP Combo
    - Interchangable judgement type (PERFECT/GREAT/GOOD/etc.)
    - Interchangable LIFE value (v3 only)

## Terms of Use

**(REQUIRED)** In the description of your video, please copy the text here:

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com) & ぴぃまん (https://pim4n-net.com)
- Original by 名無し｡ (https://sevenc7c.com) 
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

**JP**
```
プロセカ風動画作成補助ツール：
- フォーク：TootieJin様(https://tootiejin.com) & ぴぃまん様(https://pim4n-net.com)
- 作成：名無し｡様（https://sevenc7c.com）
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

> [!NOTE]
> **(optional)** You can remove watermark by check/unchecking `Watermark` in the `Root@pjsekai-overlay-en` element.\
> [![image](https://github.com/user-attachments/assets/48a636da-8ec0-443b-9cf4-b73fd93c47df)]()

## Requirements

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) ([JP Installation Guide](https://aviutl.info/dl-innsuto-ru/))
  - (**Recommended:** [AviUtl JP Installer Script](https://github.com/menndouyukkuri/aviutl-installer-script) OR [AviUtl EN Extra Pack](https://www.videohelp.com/download/AviUtl_setup_1.14.exe))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- **Fonts:** RodinNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- Basic knowledge of AviUtl

*- Refer to this [English guide](https://forum.videohelp.com/threads/360967-AviUtl-video-tutorial#post2288628) for using AviUtl EN, and [setup guide (partially outdated)](https://github.com/Khronophobia/pjsekai-overlay-english/wiki) for pjsekai-overlay.*
> [!IMPORTANT]
> **REMEMBER TO GO TO `File > SETTINGS > SYSTEM` AND SET THE `Max resolution` TO 4000x4000 (or bigger)!!!!!!!!**

## Video Guide

1. [Make your chart first.](https://wiki.purplepalette.net/create-charts)
2. Go to [Sonolus](https://sonolus.com/) to find your chart.
3. Screen record the video with **BLACK Background**, **「Stage」OFF** and **「Hide UI」ON**
4. Transfer the video file to your computer.
   - Download the [ffmpeg](https://www.ffmpeg.org/) encoder if you haven't.
5. Use the FFmpeg encoder (`ffmpeg -i source.mp4 output.mp4`)
   - **This step is required so that the video file doesn't shift speed.**
6. Once done, refer to the usage guide below.

## Usage Guide
0. Install AviUtl (**Recommended:** [AviUtl JP Installer Script](https://github.com/menndouyukkuri/aviutl-installer-script) OR [AviUtl EN Extra Pack](https://www.videohelp.com/download/AviUtl_setup_1.14.exe))
1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
2. Unzip the file
3. Go to `aviutl.exe` file location and make a new folder `script` in this directory
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *If a folder is missing, make a new folder with said name.*
4. Open AviUtl
   - **Note: You must open AviUtl before opening pjsekai-overlay-APPEND to install objects.**
5. Open `pjsekai-overlay-APPEND.exe`
   - If followed correctly, you'll be able to see 2 green lines in your console:
     ```
     AviUtlオブジェクトのインストールに成功しました。(AviUtl object successfully installed.)
     AviUtl依存関係スクリプトのインストールに成功しました。(AviUtl dependency scripts successfully installed.)
     ```
6. Input the chart ID including the prefix.
   - `chcy-`: Chart Cyanvas
      - `0`: Archive ([cc.sevenc7c.com](https://cc.sevenc7c.com))
      - `1`: Offshoot server ([chart-cyanvas.com](https://chart-cyanvas.com))
      - `Others (URL domain)`: Different Cyanvas instance
   - `ptlv-`: Potato Leaves ([ptlv.sevenc7c.com](https://ptlv.sevenc7c.com))
   - ~~`utsk-`: Untitled Sekai ([us.pim4n-net.com](https://us.pim4n-net.com))~~
   - `UnCh-`: UntitledCharts ([untitledcharts.com](https://untitledcharts.com))
   - `coconut-next-sekai-`: Next SEKAI ([coconut.sonolus.com/next-sekai](https://coconut.sonolus.com/next-sekai))
   - `lalo-`: [laoloser](https://www.youtube.com/@laoloserr)'s server ([sonolus.laoloser.com](https://laoloser.com/sonolus))
   - `sync-`: Local Server ([ScoreSync](https://github.com/Piliman22/ScoreSync))
7. Import specified exo file by navigating to your `pjsekai-overlay/dist/[chart ID]` directory:
   - **For phone users:** `main_en_16-9_1920x1080.exo`
   - **For iPad users:** `main_en_4-3_1440x1080.exo`
8. Once finished, export your video as mp4: `File > Export with plugins > Adv. x264 Export(GUI) Ex`
   - For easier use, you can download [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip) and import the file to the `Plugins` folder
---------------------------------------------------------------------------------------

# pjsekai-overlay-APPEND / フォークプロセカ風動画作成補助ツール (日本語)

[TootieJin](https://tootiejin.com)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)用フォーク。
pjsekai-overlay(-APPEND) は、プロセカの創作譜面をプロセカ風の動画にするためのオープンソースのツールです。

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - [追加アセット](./assets/extra%20assets/) ([ReiyuN](hhttps://github.com/Reiyunkun)さん、[Gaven](https://github.com/gaven1880)さんと[YumYummity](https://github.com/YumYummity)さん、ご寄稿ありがとうございました。)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成
  - iPad（4:3）動画対応
  - 英語版AviUtlの使用機能
  - v1 UIスキン（フル対応）
  - 譜面のタグ（またはタイトル）に基づいて、AviUtlで生成される譜面の難易度を自動的に変更する
  - 追加サーバーのサポート
  - 各種UIカスタマイズ
    - アニメーション付きスコア表示
    - 各要素のアニメーションON/OFF切り替え
    - 交換可能なAPコンボ
    - 交換可能な判定タイプ（PERFECT/GREAT/GOODなど）
    - 交換可能なライフ値（v3のみ）

## 利用規約

**(必須)** 動画の説明文に、こちらのテキストをコピーしてください：

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com) & ぴぃまん (https://pim4n-net.com)
- Original by 名無し｡ (https://sevenc7c.com) 
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

**JP**
```
プロセカ風動画作成補助ツール：
- フォーク：TootieJin様(https://tootiejin.com) & ぴぃまん様(https://pim4n-net.com)
- 作成：名無し｡様（https://sevenc7c.com）
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

> [!NOTE]
> **(任意)** `設定@pjsekai-overlay`要素でチェック/チェックを外すことで、`透かし`を消すことができます。\
> [![image](https://github.com/user-attachments/assets/05cc6e7d-0e62-4729-a4df-c8d8634b0a10)]()

## 必須事項

- [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest)（[導入方法](https://aviutl.info/dl-innsuto-ru/)）
  - (**強く推奨：** [AviUtl インストーラースクリプト](https://github.com/menndouyukkuri/aviutl-installer-script))
- [Unmult](https://github.com/mes51/AVIUtl_Unmult)
- **フォント：** ロダンNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- AviUtlの基本的な知識

> [!IMPORTANT]
> **必ず`ファイル > 設定 > システム`に移動し、`最大画像サイズ`を4000x4000（またはそれ以上）に設定してください!!!!!!!!**

## 動画の作り方

1. [譜面を作る](https://wiki.purplepalette.net/create-charts)
2. [Sonolus](https://sonolus.com/)で譜面を撮影する
3. **背景を黒**にし、**「Stage」OFF**と **「Hide UI」ON**にして、動画をスクリーン録画します。
4. 撮影したプレイ動画のファイルをパソコンに転送する
   - まだダウンロードしていない場合は、[ffmpeg](https://www.ffmpeg.org/) エンコーダーをダウンロードする
5. FFmpegエンコーダーを使用する（`ffmpeg -i source.mp4 output.mp4`）
   - **この手順は、動画ファイルの速度がずれないようにするために必要です。**
6. 下の利用方法に従って UI を後付けする

## 利用方法
0. AviUtlをインストールする(**強く推奨：** [AviUtl インストーラースクリプト](https://github.com/menndouyukkuri/aviutl-installer-script))
1. 右の Releases から最新のバージョンの zip を[ダウンロード](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/)する
2. zip を解凍する
3. `aviutl.exe`ファイルの場所へ移動し、ここに新しいフォルダ`script`を作成する
```
   aviutl
      ⌞Plugins
         ⌞script
```
   - *フォルダがない場合は、その名前で新しいフォルダを作ってください。*
4. AviUtl を起動する
   - **pjsekai-overlay が起動する前に AviUtl を起動するとオブジェクトのインストールが行われます。**
5. `pjsekai-overlay-APPEND.exe` を起動する
   - 正しく実行すれば、コンソールに2本の緑色の線が表示されます：
     ```
     AviUtlオブジェクトのインストールに成功しました。(AviUtl object successfully installed.)
     AviUtl依存関係スクリプトのインストールに成功しました。(AviUtl dependency scripts successfully installed.)
     ```
6. 譜面IDを接頭辞込みで入力して下さい
   - `chcy-`: Chart Cyanvas
      - `0`: アーカイブ([cc.sevenc7c.com](https://cc.sevenc7c.com))
      - `1`: 分岐サーバー([chart-cyanvas.com](https://chart-cyanvas.com))
      - `その他(URLドメイン)`: 異なるCyanvasインスタンス
   - `ptlv-`: Potato Leaves ([ptlv.sevenc7c.com](https://ptlv.sevenc7c.com))
   - ~~`utsk-`: Untitled Sekai ([us.pim4n-net.com](https://us.pim4n-net.com))~~
   - `UnCh-`: UntitledCharts ([untitledcharts.com](https://untitledcharts.com))
   - `coconut-next-sekai-`: Next SEKAI ([coconut.sonolus.com/next-sekai](https://coconut.sonolus.com/next-sekai))
   - `lalo-`: [laoloser](https://www.youtube.com/@laoloserr)のサーバー([sonolus.laoloser.com](https://laoloser.com/sonolus))
   - `sync-`: ローカルサーバー([ScoreSync](https://github.com/Piliman22/ScoreSync))
7. `pjsekai-overlay/dist/[譜面ID]`ディレクトリに移動して、指定したexoファイルをインポートします：
   - **スマホ向け:** `main_jp_16-9_1920x1080.exo`
   - **iPad向け:** `main_jp_4-3_1440x1080.exo`
8. 完了後、動画をmp4形式で出力します：`ファイル > プラグイン出力 > 拡張 x264 出力(GUI) Ex`
   - 使いやすくするため、[easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip)をダウンロードし、ファイルを`Plugins`フォルダにインポートできます
