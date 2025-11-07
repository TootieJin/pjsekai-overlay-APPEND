[**English Section**](#pjsekai-overlay-append--forked-pjsekai-style-video-creation-tool) | [**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール)

[![Releases](https://img.shields.io/github/downloads/TootieJin/pjsekai-overlay-APPEND/total?style=plastic)](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/) [![Stargazers](https://img.shields.io/github/stars/TootieJin/pjsekai-overlay-APPEND?style=plastic&color=yellow)](https://github.com/TootieJin/pjsekai-overlay-APPEND/stargazers)
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool

[![pjsekai-overlay-APPEND thumbnail](https://github.com/user-attachments/assets/dcc037f7-1c2b-4b83-b17b-b3c5155d670c)]()

Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com), an open-sourced tool to make Project Sekai Fanmade (custom chart) videos.

> [!CAUTION]
> **For English users:** This tool is primary only for people with technical know-how and basic knowledge of AviUtl.\
> Only use this tool if you can figure it out yourself. **DO NOT open issues, DM me, or request help about this**.
> 
> *Also, to [a certain someone](https://discordid.netlify.app/?id=1370076899404939327) (a.k.a [this person](https://discordid.netlify.app/?id=919036186473947187)) with the mindset of `"Just switch to a different editing software since I don’t even know how to install aviutl"` [(image source)](https://github.com/user-attachments/assets/4850442d-3f3a-438a-92d9-97d052f2fba0): I suggest you **make your own pjsekai-overlay that supports your desired editing software.** (if you can even find a video editor that is as versatile and extensible as AviUtl, that is).*

- **16:9**

https://github.com/user-attachments/assets/dda7225a-a7f3-41d4-bbf4-9cec9b03b840

- **4:3 (Tournament Mode ON)**

https://github.com/user-attachments/assets/ab4ee52c-2ffa-4941-b916-87e1f3559d72

- **v1 Skin (1e+30 power)**

https://github.com/user-attachments/assets/3efab743-246a-4da7-8d80-a02b2f09f5b3

- **AviUtl ExEdit2 preview**

<img width="2560" height="1554" alt="image" src="https://github.com/user-attachments/assets/18f6d16b-5ba4-4953-aa7d-5ceebb87348a" />

- **Video Example**

*(Click the image to watch it)*\
[![【Project Sekai x Honkai: Star Rail】Nameless Faces - HOYO-MiX feat. Lilas Ikuta (Fanmade)](https://img.youtube.com/vi/uXx1OZDQZOI/maxresdefault.jpg)](https://youtu.be/uXx1OZDQZOI)
[![【Project Sekai Fanmade? (v3→v1)】Hello, SEKAI - DECO*27【ETERNAL Lv32】](https://img.youtube.com/vi/BHVNuwxA1ek/maxresdefault.jpg)](https://youtu.be/BHVNuwxA1ek)

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - **AviUtl ExEdit2 Full Support**
  - [Extra assets](./extra%20assets)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos
  - iPad (4:3) video support
  - Ability to use the English AviUtl
  - v1 UI skin (Full support)
  - Automatically changes chart difficulty to generate in AviUtl based on chart tag (or title)
  - Additional support for more servers
  - **[Various UI customization](#ui-customization-specifications)**
    - Animated Scoring
    - Adjust animation speed in different elements
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
- フォーク：TootieJin氏 (https://tootiejin.com) & ぴぃまん氏 (https://pim4n-net.com)
- 作成：名無し｡氏 (https://sevenc7c.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

> [!NOTE]
> **(optional)** You can remove watermark by check/unchecking `Watermark` in the `Root@pjsekai-overlay-en` element.
> 
> <img width=80% height=80% alt="image" src="https://github.com/user-attachments/assets/48a636da-8ec0-443b-9cf4-b73fd93c47df" />

## Requirements

| AviUtl                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | **AviUtl ExEdit2 (Beta)**                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip) ([JP Installation Guide](https://aviutl.info/dl-innsuto-ru/))<br>- **JP version:** [AviUtl JP Installer Script](https://github.com/menndouyukkuri/aviutl-installer-script)<br>- **EN version:** [AviUtl EN Extra Pack](https://www.videohelp.com/download/AviUtl_setup_1.14.exe)<br><br>*Import all plugins to `aviutl\Plugin` folder* | [AviUtl ExEdit2 **(via installer)**](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [MP4Exporter](https://apps.esugo.net/aviutl2-mp4exporter/)<br>- English language can be selected in `設定 > 言語の設定 > English (Settings > Language > English)`<br><br>*Import `AviUtl2\lwinput.aui2` & `MP4Exporter.auo2` plugin to `C:\ProgramData\aviutl2\Plugin` folder* |
- **Fonts:** RodinNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- Basic knowledge of AviUtl

## Video Guide

0. Make your chart first.
1. Go to [Sonolus](https://sonolus.com/) to find your chart.
2. Screen record the video with **BLACK Background**, **「Stage」OFF** and **「Hide UI」ON**
3. Transfer the video file to your computer.
4. Download the [FFmpeg](https://www.ffmpeg.org/) encoder and use it (`ffmpeg -i source.mp4 output.mp4`)
   - **This step is required so that the video file doesn't shift speed. (not known to be needed in AviUtl ExEdit2 yet)**
   - Alternatively, you can use your video editor to export it.
5. Once done, refer to the usage guide below.

## Usage Guide
0. [Install AviUtl](#requirements)
1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
   - **(optional)** You can download `extra-assets.zip` for use of non-essential assets.
2. Unzip the file in a directory that **doesn't contain non-ASCII characters nor require administrative permissions** (see https://github.com/TootieJin/pjsekai-overlay-APPEND/issues/5)
3. Follow the steps based on which AviUtl you want to use:

| AviUtl                                                                                                                                                                                                                                                                                                                           	| **AviUtl ExEdit2 (Beta)**                                                                      	|
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------	|
| 1. Go to `aviutl.exe` file location and make a new folder script in the `aviutl\Plugins\script` directory<br>2. **Install UI objects** (open AviUtl > open `pjsekai-overlay-APPEND.exe`)<br>- **Note: You must open AviUtl before opening pjsekai-overlay-APPEND to install objects**<br>3. **Restart AviUtl to apply changes.** 	| 1. Open `pjsekai-overlay-APPEND.exe`<br>2. Press `2` to choose the AviUtl ExEdit2 instance 	|

4. In your `pjsekai-overlay-APPEND.exe` command prompt, input the chart ID including the prefix.
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
5. Import specified exo/alias(.object) file by navigating to your `pjsekai-overlay/dist/[chart ID]` directory:

|      	| AviUtl                                                          	| **AviUtl ExEdit2 (Beta)**                                        	|
|------	|-----------------------------------------------------------------	|------------------------------------------------------------------	|
| 16:9 	| `main_en_16-9_1920x1080.exo`<br>`v1-skin_en_16-9_1920x1080.exo` 	| `main2_16-9_1920x1080.object`<br>`v1-skin_16-9_1920x1080.object` 	|
| 4:3  	| `main_en_4-3_1440x1080.exo`<br>`v1-skin_en_4-3_1440x1080.exo`   	| `main2_4-3_1440x1080.object`<br>`v1-skin_4-3_1440x1080.object`   	|

6. Once finished, export your video as mp4

## UI Customization Specifications
> [!TIP]
> You can change default values in `default.ini`, **reducing the need to change values in AviUtl.**
> 
> <img width=70% height=70% alt="image" src="https://github.com/user-attachments/assets/bd2151ad-2045-46ef-96c2-cbbb426a68d9" />

### Root@pjsekai-overlay-en
<img width="383" height="79" alt="image" src="https://github.com/user-attachments/assets/76da6c77-f6a7-4480-b279-b5d53f3e583f" />

| **Name**      	| Description                                                                                                	| Default 	|        Range       	|
|---------------	|------------------------------------------------------------------------------------------------------------	|:-------:	|:------------------:	|
| **Offset**<br>`offset`    	| Number of frames to shift events<br>- Increase to shift timing later<br>- Decrease to shift timing earlier 	|  216.0  	| -99999.9 ~ 99999.9 	|
| **Cache**<br>`cache`     	| When cache is set to 0, any change in the `data.ped` is applied immediately                               	|    1    	|        0 ~ 1       	|
| **Font type**<br>`font_type` 	| Set font type configuration for the watermark text<br>(`0` - メイリオ, `1` - RodinNTLG EB)                 	|    0    	|        0 ~ 1       	|
| **Watermark**<br>`watermark` 	| Enable watermark text at the bottom-left corner                                                            	|   1(ON)  	|          X         	|

### Life@pjsekai-overlay-en (v3 only)
<img width="125" height="125" alt="LifeUP" src="https://github.com/user-attachments/assets/6f7a7db8-50bb-43cf-9463-5f46325c862e" /> <img width=50% height=50% alt="life" src="https://github.com/user-attachments/assets/7aab3534-66cf-4dad-936e-3d423ecce615" />

| **Name** 	| Description                                                                         	| Default 	|   Range  	|
|----------	|-------------------------------------------------------------------------------------	|:-------:	|:--------:	|
| **LIFE**<br>`life` 	| LIFE value (self-explanatory)<br>- When value changes, the LIFE bar changes as well 	|   1000  	| 0 ~ 9999 	|
| **Show overflow LIFE bar**<br>`overflow` 	| <img width=50% height=50% alt="life_overflow" src="https://github.com/user-attachments/assets/75ad981f-cb1d-4112-939d-8f8bf39a1222" /> 	|  0(OFF)  	|    X   	|

### Score@pjsekai-overlay-en
<img width="125" height="125" alt="ScoreUP" src="https://github.com/user-attachments/assets/a5a8b0f0-035c-4951-8ae3-d2038945d86c" /> <img width=50% height=50% alt="bg" src="https://github.com/user-attachments/assets/3db93b3e-2280-46e1-a08f-00e50a5e5e8c" />

| **Name**             	| Description                                    	| Default 	|  Range 	|
|----------------------	|------------------------------------------------	|:-------:	|:------:	|
| **Min Digit**<br>`min_digit`        	| Render the minimum amount of digits in score   	|    8    	| 1 ~ 17 	|
| **Animated Scoring**<br>`anim_score` 	| Increase incrementally rather than all at once 	|  0(OFF)  	|    X   	|
| **Animation Speed**<br>`speed`   | Adjust animation speed                       	   |   1.00     |    X 	   |

### Combo@pjsekai-overlay-en
<img width="148" height="49" alt="pt" src="https://github.com/user-attachments/assets/9db50558-cf81-4ed8-a2bd-1d4bbd22e156" /> <img width="145" height="45" alt="nt" src="https://github.com/user-attachments/assets/3ca0f65e-8ce6-40c9-8ff4-53a8ee9d2f81" />

| **Name**                          	| Description                            	  | Default | Range 	|
|-----------------------------------	|----------------------------------------	  |:-------:|:-----:	|
| **AP Combo**<br>`ap`                      	| Toggle AP Combo status                 	  |    1   	| 0 ~ 1 	|
| **Combo Tag**<br>`tag`                     	| Toggle rendering combo tag             	  |    1   	| 0 ~ 1 	|
| **Animation Speed**<br>`speed`                 	| Adjust animation speed               	  |   1.00   |   X 	|
| **Show the last X digit(s)**<br>`last_digit`        	| `4`: 12345 -> /2345<br>`2`: 12345 -> ///45 |    4   	|   X   	|
| **n00 Combo Burst**<br>`combo_burst`        	| Toggle World Dai Star's combo burst effect |    0(OFF)   	|   X   	|

### Judgement@pjsekai-overlay-en
<img width="125" height="125" alt="SkillUP" src="https://github.com/user-attachments/assets/e29f426d-71ae-4de5-912a-a5c7375f538d" />

| **Name**       	| Description                                                                   	| Default 	| Range 	|
|----------------	|-------------------------------------------------------------------------------	|:-------:	|:-----:	|
| **Judge Type**<br>`judge` 	| `1`: <img width=25% height=25% alt="perfect" src="https://github.com/user-attachments/assets/28950e9e-0dac-49d9-81d3-70bdaa2d6f0c" /><br>`2`: <img width=25% height=25% alt="great" src="https://github.com/user-attachments/assets/ccf333a7-795d-43ad-8002-a9d2220e18a6" /><br>`3`: <img width=25% height=25% alt="good" src="https://github.com/user-attachments/assets/9d0a26bb-c8e7-47d0-9a3d-717b4ad0e0fa" /><br>`4`: <img width=25% height=25% alt="bad" src="https://github.com/user-attachments/assets/5b757195-8bd4-4beb-9f77-808000f1d865" /><br>`5`: <img width=25% height=25% alt="miss" src="https://github.com/user-attachments/assets/734ead15-491b-4bdb-9017-f2b30ab32223" /><br>`6`: <img width=25% height=25% alt="auto" src="https://github.com/user-attachments/assets/b9d674cf-1b69-478e-b2be-53691109b12d" /> 	|    1    	| 1 ~ 6 	|
| **Animation Speed**<br>`speed`                 	| Adjust animation speed               	  |   1.00   |   X 	|

## Special Thanks
- **[@sevenc-nanashi](https://github.com/sevenc-nanashi) for developing the [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) & [pjsekai-background-gen-rust](https://github.com/sevenc-nanashi/pjsekai-background-gen-rust) tool!**
- [@Piliman22](https://github.com/Piliman22) for contribution to make [ScoreSync](https://github.com/Piliman22/ScoreSync) in pjsekai-overlay possible!
- [@Reiyunkun](https://github.com/Reiyunkun), [@gaven1880](https://github.com/gaven1880), and [@YumYummity](https://github.com/YumYummity) for providing [additional PJSK assets](./extra%20assets)!
- [@Khronophobia](https://github.com/Khronophobia) for the customized lane assets in Blender!
- [@nagotown](https://github.com/nagotown) for the [aviutl2-EN](https://github.com/nagotown/aviutl2-EN) translation script!
- And everyone who used my tool... thank you all so much.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

# pjsekai-overlay-APPEND / フォークプロセカ風動画作成補助ツール

[TootieJin](https://tootiejin.com)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)用フォーク。pjsekai-overlay(-APPEND) は、プロセカの創作譜面をプロセカ風の動画にするためのオープンソースのツールです。

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - **AviUtl ExEdit2 フル対応**
  - [追加アセット](./extra%20assets/)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成
  - iPad（4:3）動画対応
  - 英語版AviUtlの使用機能
  - v1 UIスキン（フル対応）
  - 譜面のタグ（またはタイトル）に基づいて、AviUtlで生成される譜面の難易度を自動的に変更する
  - 追加サーバーのサポート
  - **[各種UIカスタマイズ](#uiカスタマイズ仕様書)**
    - アニメーション付きスコア表示
    - 異なる要素でアニメーション速度を調整する
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
- フォーク：TootieJin氏 (https://tootiejin.com) & ぴぃまん氏 (https://pim4n-net.com)
- 作成：名無し｡氏 (https://sevenc7c.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

> [!NOTE]
> **(任意)** `設定@pjsekai-overlay`要素でチェック/チェックを外すことで、`透かし`を消すことができます。
> 
> <img width=80% height=80% alt="image" src="https://github.com/user-attachments/assets/05cc6e7d-0e62-4729-a4df-c8d8634b0a10" />

## 必須事項
| AviUtl                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | **AviUtl ExEdit2**                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip)<br>- **おすすめ：** [AviUtl インストーラースクリプト](https://github.com/menndouyukkuri/aviutl-installer-script)<br><br>*すべてのプラグインを `aviutl\Plugin` フォルダにインポートする* | [AviUtl ExEdit2 **(インストーラ版)**](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [MP4Exporter](https://apps.esugo.net/aviutl2-mp4exporter/)<br><br>*`AviUtl2\lwinput.aui2`と`MP4Exporter.auo2`プラグインを`C:\ProgramData\aviutl2\Plugin`フォルダにインポートする* |
- **フォント：** ロダンNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- AviUtlの基本的な知識

## 動画の作り方

0. 譜面を作る
1. [Sonolus](https://sonolus.com/)で譜面を撮影する
2. **背景を黒**にし、**「Stage」OFF** と **「Hide UI」ON**にして、動画をスクリーン録画します。
3. 撮影したプレイ動画のファイルをパソコンに転送する
   - まだダウンロードしていない場合は、[FFmpeg](https://www.ffmpeg.org/) エンコーダーをダウンロードする
4. [FFmpeg](https://www.ffmpeg.org/) エンコーダーをダウンロードして使用します（`ffmpeg -i source.mp4 output.mp4`）。
   - **この手順は、動画ファイルの速度がずれないようにするために必要です。(AviUtl ExEdit2では現時点で必要とされていない)**
   - 代わりに、動画編集ソフトでエクスポートすることもできます。
5. 下の利用方法に従って UI を後付けする

## 利用方法
0. [AviUtlをインストールする](#必須事項)
1. 右のReleasesから最新のバージョンのzipを[ダウンロード](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/)する
   - **(任意)** 必須ではないアセットを使用する場合は、`extra-assets.zip` をダウンロードできます。
2. ファイルを解凍するディレクトリは、**非ASCII文字を含まず、管理者権限を必要としない**ものを選んでください。（詳細は https://github.com/TootieJin/pjsekai-overlay-APPEND/issues/5 を参照）
3. 使用したいAviUtlに応じて以下の手順に従ってください：

| AviUtl                                                                                                                                                                                                                                                                                                                           	| **AviUtl ExEdit2 (Beta)**                                                                      	|
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------	|
| 1. `aviutl.exe` ファイルの場所へ移動し、`aviutl\Plugins\script` ディレクトリ内に新しいフォルダ「script」を作成してください<br>2. **UIオブジェクトのインストール** (AviUtlを起動 > `pjsekai-overlay-APPEND.exe` を起動)<br>- **注意: オブジェクトをインストールするには、pjsekai-overlay-APPENDを開く前に必ずAviUtlを起動しておく必要があります**<br>3. **変更を適用するためAviUtlを再起動してください。**     | 1. `pjsekai-overlay-APPEND.exe`を開く<br>2. `2`を押してAviUtl ExEdit2インスタンスを選択     |

4. `pjsekai-overlay-APPEND.exe`コンソールで、譜面IDを接頭辞込みで入力して下さい。
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
5. `pjsekai-overlay/dist/[譜面ID]`ディレクトリに移動して、指定したexo/alias(.object)ファイルをインポートします：

|      	| AviUtl                                                          	| **AviUtl ExEdit2 (Beta)**                                        	|
|------	|-----------------------------------------------------------------	|------------------------------------------------------------------	|
| 16:9 	| `main_jp_16-9_1920x1080.exo`<br>`v1-skin_jp_16-9_1920x1080.exo` 	| `main2_16-9_1920x1080.object`<br>`v1-skin_16-9_1920x1080.object` 	|
| 4:3  	| `main_jp_4-3_1440x1080.exo`<br>`v1-skin_jp_4-3_1440x1080.exo`   	| `main2_4-3_1440x1080.object`<br>`v1-skin_4-3_1440x1080.object`   	|

6. 完了後、動画をmp4形式で出力します

## UIカスタマイズ仕様書
> [!TIP]
> `default.ini`でデフォルト値を変更できます。**これによりAviUtlで値を変更する必要が軽減されます。**
> 
> <img width=70% height=70% alt="image" src="https://github.com/user-attachments/assets/bd2151ad-2045-46ef-96c2-cbbb426a68d9" />

### 設定@pjsekai-overlay
<img width="383" height="79" alt="image" src="https://github.com/user-attachments/assets/76da6c77-f6a7-4480-b279-b5d53f3e583f" />

| **名前**         | 説明                                                                                                    | デフォルト     |        範囲           |
|---------------   |------------------------------------------------------------------------------------------------------------	|:-------:	|:------------------:	|
| **オフセット**<br>`offset`    | タイミングをシフトするフレーム数<br>- 増加するとタイミングが遅くなる<br>- 減少するとタイミングが早くなる 	|  216.0  	| -99999.9 ~ 99999.9 	|
| **キャッシュ**<br>`cache`    	| キャッシュが0に設定されている場合、`data.ped`の変更は即時反映されます                               	|    1    	|        0 ~ 1       	|
| **フォント種類**<br>`font_type`  | 透かしテキストのフォント種類設定を設定する<br>(`0` - メイリオ, `1` - RodinNTLG EB)                 	|    0    	|        0 ~ 1       	|
| **透かし**<br>`watermark` 	   | 左下隅に透かしテキストを表示する                                                            	|   1(ON)  	|          X         	|

### ライフ@pjsekai-overlay (v3のみ)
<img width="125" height="125" alt="LifeUP" src="https://github.com/user-attachments/assets/6f7a7db8-50bb-43cf-9463-5f46325c862e" /> <img width=50% height=50% alt="life" src="https://github.com/user-attachments/assets/7aab3534-66cf-4dad-936e-3d423ecce615" />

| **名前**         | 説明                                                                                                    | デフォルト     |        範囲           |
|----------	|-------------------------------------------------------------------------------------	|:-------:	|:--------:	|
| **ライフ**<br>`life` 	| LIFEの値（自明）<br>- 値が変化すると、LIFEバーも変化します 	|   1000  	| 0 ~ 9999 	|
| **過剰なライフバー**<br>`overflow` 	| <img width=50% height=50% alt="life_overflow" src="https://github.com/user-attachments/assets/75ad981f-cb1d-4112-939d-8f8bf39a1222" /> 	|  0(OFF)  	|    X   	|

### スコア@pjsekai-overlay
<img width="125" height="125" alt="ScoreUP" src="https://github.com/user-attachments/assets/a5a8b0f0-035c-4951-8ae3-d2038945d86c" /> <img width=50% height=50% alt="bg" src="https://github.com/user-attachments/assets/3db93b3e-2280-46e1-a08f-00e50a5e5e8c" />

| **名前**         | 説明                                                                                                    | デフォルト     |        範囲           |
|----------------------	|------------------------------------------------	|:-------:	|:------:	|
| **最小桁数**<br>`min_digit`        	| スコアの桁数を最小限に表示する   	|    8    	| 1 ~ 17 	|
| **アニメーション採点**<br>`anim_score` 	| 一気にではなく、段階的に増やす 	|  0(OFF)  	|    X   	|
| **アニメーション速度**<br>`speed`	| アニメーション速度を調整する 	|    1.00    	|   X 	|

### コンボ@pjsekai-overlay
<img width="148" height="49" alt="pt" src="https://github.com/user-attachments/assets/9db50558-cf81-4ed8-a2bd-1d4bbd22e156" /> <img width="145" height="45" alt="nt" src="https://github.com/user-attachments/assets/3ca0f65e-8ce6-40c9-8ff4-53a8ee9d2f81" />

| **名前**         | 説明                                                                                                    | デフォルト     |        範囲           |
|-----------------------------------	|----------------------------------------	|:-------:	|:-----:	|
| **APコンボ**<br>`ap`                      	| APコンボ状態を切り替える                 	|    1    	| 0 ~ 1 	|
| **コンボタグ**<br>`tag`                     	| コンボタグの表示切り替え             	|    1    	| 0 ~ 1 	|
| **アニメーション速度**<br>`speed`                     	| アニメーション速度を調整する 	|    1.00    	| X 	|
| **最後のX桁を表示**<br>`last_digit` 	| `4`: 12345 -> /2345<br>`2`: 12345 -> ///45 |    4   	|   X   	|
| **n00コンボ効果**<br>`combo_burst`            | ユメステのn00コンボ効果を切り替え |    0(OFF)       |   X       |

### 判定@pjsekai-overlay
<img width="125" height="125" alt="SkillUP" src="https://github.com/user-attachments/assets/e29f426d-71ae-4de5-912a-a5c7375f538d" />

| **名前**         | 説明                                                                                                    | デフォルト     |        範囲           |
|----------------	|-------------------------------------------------------------------------------	|:-------:	|:-----:	|
| **判定タイプ**<br>`judge` 	| `1`: <img width=25% height=25% alt="perfect" src="https://github.com/user-attachments/assets/28950e9e-0dac-49d9-81d3-70bdaa2d6f0c" /><br>`2`: <img width=25% height=25% alt="great" src="https://github.com/user-attachments/assets/ccf333a7-795d-43ad-8002-a9d2220e18a6" /><br>`3`: <img width=25% height=25% alt="good" src="https://github.com/user-attachments/assets/9d0a26bb-c8e7-47d0-9a3d-717b4ad0e0fa" /><br>`4`: <img width=25% height=25% alt="bad" src="https://github.com/user-attachments/assets/5b757195-8bd4-4beb-9f77-808000f1d865" /><br>`5`: <img width=25% height=25% alt="miss" src="https://github.com/user-attachments/assets/734ead15-491b-4bdb-9017-f2b30ab32223" /><br>`6`: <img width=25% height=25% alt="auto" src="https://github.com/user-attachments/assets/b9d674cf-1b69-478e-b2be-53691109b12d" /> 	|    1    	| 1 ~ 6 	|
| **アニメーション速度**<br>`speed`                     	| アニメーション速度を調整する 	|    1.00    	| X 	|

## 特別なお礼
- **[@sevenc-nanashi](https://github.com/sevenc-nanashi)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)と[pjsekai-background-gen-rust](https://github.com/sevenc-nanashi/pjsekai-background-gen-rust)ツールの開発に感謝します！**
- [@Piliman22](https://github.com/Piliman22)氏、pjsekai-overlayにおける[ScoreSync](https://github.com/Piliman22/ScoreSync)の実現に貢献いただき感謝申し上げます！
- [@Reiyunkun](https://github.com/Reiyunkun)氏、[@gaven1880](https://github.com/gaven1880)氏、[@YumYummity](https://github.com/YumYummity)氏による[追加のプロセカアセット](./extra%20assets)の提供に感謝します！
- [@Khronophobia](https://github.com/Khronophobia)氏にはBlenderでのカスタムレーンアセットを提供いただきました！
- [@nagotown](https://github.com/nagotown)氏による [aviutl2-EN](https://github.com/nagotown/aviutl2-EN) 翻訳スクリプト！
- そして私のツールを使ってくださった皆様、本当にありがとうございます。
