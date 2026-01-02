[**English Section**](#pjsekai-overlay-append--forked-pjsekai-style-video-creation-tool) | [**日本語セクション**](#pjsekai-overlay-append--フォークプロセカ風動画作成補助ツール)

[![Releases](https://img.shields.io/github/downloads/TootieJin/pjsekai-overlay-APPEND/total?style=plastic)](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/) [![Stargazers](https://img.shields.io/github/stars/TootieJin/pjsekai-overlay-APPEND?style=plastic&color=yellow)](https://github.com/TootieJin/pjsekai-overlay-APPEND/stargazers)
# pjsekai-overlay-APPEND / Forked PJSekai-style video creation tool

[![pjsekai-overlay-APPEND thumbnail](https://github.com/user-attachments/assets/dcc037f7-1c2b-4b83-b17b-b3c5155d670c)]()

Fork of [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) by [TootieJin](https://tootiejin.com), an open-sourced tool to make Project Sekai Fanmade (custom chart) videos - in other words... "Make your Sonolus look like Project Sekai."

> [!CAUTION]
> **Before using pjsekai-overlay-APPEND:** This tool is primary only for people with **technical know-how and basic knowledge of AviUtl.** \
> Only use this tool if you can figure it out yourself. If problems arise, **please [make a discussion thread.](https://github.com/TootieJin/pjsekai-overlay-APPEND/discussions)** \
> Please also read the [Terms of Use](#terms-of-use) after using the tool.
> 
> - *P.S. To [a certain someone](https://discordid.netlify.app/?id=1370076899404939327) (a.k.a [this person](https://discordid.netlify.app/?id=919036186473947187)) with the mindset of `"Just switch to a different editing software since I don’t even know how to install aviutl"` [(image source)](https://github.com/user-attachments/assets/4850442d-3f3a-438a-92d9-97d052f2fba0): I suggest you **make your own pjsekai-overlay that supports your desired editing software.** (if you can even find a video editor that is as versatile and extensible as AviUtl, that is).*

This is a forked version of pjsekai-overlay with additional features originally not in the main repo, including:
  - **AviUtl ExEdit2 Full Support**
  - [Extra assets](./extra-assets)
  - Added/adjusted elements to look identical to the official photography
  - Quickly make 1080p videos
  - iPad (4:3) video support
  - Ability to use the English AviUtl
  - v1 UI skin (Full support)
  - Automatically changes chart difficulty to generate in AviUtl based on chart tag (or title)
  - **(AviUtl ExEdit2 ONLY)** Supports very big number (~ 1.7e+307) using [lua-bignumber](https://github.com/thenumbernine/lua-bignumber)
  - **[Various UI customization](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki/Start-Here!#ui-customization-specifications)**
    - Animated Scoring
    - Adjust animation speed in different elements
    - Interchangable AP Combo
    - Interchangable judgement type (PERFECT/GREAT/GOOD/etc.)
    - Interchangable LIFE value
    - Achievement Rate
  - **Additional support for more servers:**
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

### • 16:9

https://github.com/user-attachments/assets/dda7225a-a7f3-41d4-bbf4-9cec9b03b840

### • 4:3 (Tournament Mode ON)

https://github.com/user-attachments/assets/ab4ee52c-2ffa-4941-b916-87e1f3559d72

### • v1 Skin (1e+30 power)

https://github.com/user-attachments/assets/3efab743-246a-4da7-8d80-a02b2f09f5b3

### • AviUtl ExEdit2 Preview

<img width=70% height=70% alt="AviUtl ExEdit2 Preview" src="https://github.com/user-attachments/assets/18f6d16b-5ba4-4953-aa7d-5ceebb87348a" /><br>

## Video Example (Click image to watch)
  - TootieJin

<a href="https://youtu.be/uXx1OZDQZOI"><img width=32% height=32% alt="【Project Sekai x Honkai: Star Rail】Nameless Faces - HOYO-MiX feat. Lilas Ikuta (Fanmade)" src="https://img.youtube.com/vi/uXx1OZDQZOI/maxresdefault.jpg" /></a> <a href="https://youtu.be/BHVNuwxA1ek"><img width=32% height=32% alt="【Project Sekai Fanmade? (v3→v1)】Hello, SEKAI - DECO*27【ETERNAL Lv32】" src="https://img.youtube.com/vi/BHVNuwxA1ek/maxresdefault.jpg" /></a> <a href="https://youtu.be/O2GxrLWbwlY"><img width=32% height=32% alt="【Project Sekai Fanmade】Confessions of a Rotten Girl - SAWTOWNE (APPEND Lv36)" src="https://img.youtube.com/vi/O2GxrLWbwlY/maxresdefault.jpg" /></a>
<a href="https://youtu.be/yyoTAJx-6ok"><img width=32% height=32% alt="【Project Sekai Fanmade (?)】Anonymous M - PinocchioP (CO-OP CHART)" src="https://img.youtube.com/vi/yyoTAJx-6ok/maxresdefault.jpg" /></a>

  - Hydrogen

<a href="https://youtu.be/aXAeV8nggPs"><img width=32% height=32% alt="【創作譜面】思いついてしまったんだ...... 曲: オーバーライド - 作: 吉田夜世" src="https://img.youtube.com/vi/aXAeV8nggPs/maxresdefault.jpg" /></a> <a href="https://youtu.be/Jqf5eb8neOA"><img width=32% height=32% alt="【創作譜面】易しいギミック 曲: OCO LIMBO - 作: MYUKKE." src="https://img.youtube.com/vi/Jqf5eb8neOA/maxresdefault.jpg" /></a> <a href="https://youtu.be/zOqBCZKK5IY"><img width=32% height=32% alt="【創作譜面】久々の公式レギュ 曲: カルチャー - 作: キタニタツヤ" src="https://img.youtube.com/vi/zOqBCZKK5IY/maxresdefault.jpg" /></a>

  - まてぃ

<a href="https://youtu.be/D0VD5fdQq3M"><img width=32% height=32% alt="マシュマロ / Marshmallow Master Lv33 【プロセカ創作譜面】" src="https://img.youtube.com/vi/D0VD5fdQq3M/maxresdefault.jpg" /></a> <a href="https://youtu.be/yvdWJUD1D9A"><img width=32% height=32% alt="アベリア APPEND Lv35 【プロセカ創作譜面】" src="https://img.youtube.com/vi/yvdWJUD1D9A/maxresdefault.jpg" /></a>

  - ReiyuN

<a href="https://youtu.be/R4jHeilZ3uk"><img width=32% height=32% alt="【 プロセカ 創作譜面 】wiege | Alien Stage x Project SEKAI【 APPEND 27 】" src="https://img.youtube.com/vi/R4jHeilZ3uk/maxresdefault.jpg" /></a> <a href="https://youtu.be/PTuaPc2vSyo"><img width=32% height=32% alt="【 Project SEKAI Fanmade 】Static【 MASTER 30 】" src="https://img.youtube.com/vi/PTuaPc2vSyo/maxresdefault.jpg" /></a>

  - Runakkushia

<a href="https://youtu.be/X4HoJvo6xRE"><img width=32% height=32% alt="【Project Sekai Fanmade】Don’t Believe in T feat. Hatsune Miku & Kasane Teto - PinocchioP" src="https://img.youtube.com/vi/X4HoJvo6xRE/maxresdefault.jpg" /></a>

## Terms of Use

1. **(REQUIRED)** In the description of your video, you must copy the text here:

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com) & ぴぃまん (https://pim4n-net.com)
- Developed by 名無し｡ (https://sevenc7c.com) 
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

**JP**
```
プロセカ風動画作成補助ツール：
- フォーク：TootieJin氏 (https://tootiejin.com) & ぴぃまん氏 (https://pim4n-net.com)
- 作成：名無し｡氏 (https://sevenc7c.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

2. This tool **should not be used for malicious purposes** (such as spreading misinformation on social media).
3. The author **assumes no responsibility whatsoever** for any issues or disadvantages arising from the use of this tool.

## Requirements

| AviUtl                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | **AviUtl ExEdit2 (Recommended)**                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip) ([JP Installation Guide](https://aviutl.info/dl-innsuto-ru/))<br>- **JP installer:** [AviUtl JP Installer Script](https://github.com/menndouyukkuri/aviutl-installer-script)<br>- **EN installer:** [AviUtl EN Extra Pack](https://www.videohelp.com/download/AviUtl_setup_1.14.exe)<br><br>*Import all plugins to `aviutl\Plugins` folder* | [AviUtl ExEdit2 **(via installer)**](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [MP4Exporter](https://apps.esugo.net/aviutl2-mp4exporter/)<br>- English language can be selected in `設定 > 言語の設定 > English` (`Settings > Language > English`)<br><br>*Import `AviUtl2\lwinput.aui2` & `MP4Exporter.auo2` plugin to `C:\ProgramData\aviutl2\Plugin` folder* |
- **Fonts:** RodinNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- Your system locale **must be `Japanese (Japan)`** [(Go here for how to change it)](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki/Start-Here!#changing-system-locale)
- **Japanese Language Pack** must be installed [(Go here for installation instructions)](https://support.microsoft.com/en-us/windows/language-packs-for-windows-a5094319-a92d-18de-5b53-1cfc697cfca8)
- Basic knowledge of AviUtl

## Usage Guide

**We've moved the guide to [the wiki page](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki).**

> [!NOTE]
> Please check the wiki often, as we may update information.

## Special Thanks
- **[@sevenc-nanashi](https://github.com/sevenc-nanashi) for developing the [pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay) & [pjsekai-background-gen-rust](https://github.com/sevenc-nanashi/pjsekai-background-gen-rust) tool!**
- [@Piliman22](https://github.com/Piliman22) for contribution to make [ScoreSync](https://github.com/Piliman22/ScoreSync) in pjsekai-overlay possible!
- [@Reiyunkun](https://github.com/Reiyunkun), [@gaven1880](https://github.com/gaven1880), and [@YumYummity](https://github.com/YumYummity) for providing [additional PJSK assets](./extra-assets)!
- [@Khronophobia](https://github.com/Khronophobia) for the customized lane assets in Blender!
- [@nagotown](https://github.com/nagotown) for the [aviutl2-EN](https://github.com/nagotown/aviutl2-EN) translation script!
- [@MattMayuga](https://github.com/MattMayuga) for the [customized judgement fonts](https://github.com/Tiny-Foxes/JudgeFonts-by-MattMayuga)!
- And everyone who used my tool, thank you all so much.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

# pjsekai-overlay-APPEND / フォークプロセカ風動画作成補助ツール

[TootieJin](https://tootiejin.com)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)用フォーク。pjsekai-overlay(-APPEND) は、プロセカの創作譜面をプロセカ風の動画にするためのオープンソースのツールです。

> [!CAUTION]
> **pjsekai-overlay-APPEND使用前に：** このツールは主に、**技術的な知識とAviUtlの基本的な理解がある方**のみを対象としています。\
> 自分で理解できる方のみご使用ください。問題が発生した場合は、**[議論スレッドを作成してください。](https://github.com/TootieJin/pjsekai-overlay-APPEND/discussions)** \
> ツール使用後は、[利用規約](#利用規約)も必ずお読みください。

これはpjsekai-overlayのフォーク版で、元々メインレポにはない以下のような追加機能があります：
  - **AviUtl ExEdit2 フル対応**
  - [追加アセット](./extra-assets/)
  - 本家撮影と同じように見えるように要素を追加/調整
  - 1080p動画を素早く作成
  - iPad（4:3）動画対応
  - 英語版AviUtlの使用機能
  - v1 UIスキン（フル対応）
  - 譜面のタグ（またはタイトル）に基づいて、AviUtlで生成される譜面の難易度を自動的に変更する
  - **(AviUtl ExEdit2 のみ)** [lua-bignumber](https://github.com/thenumbernine/lua-bignumber)を使用して非常に大きな数値（~ 1.7e+307）を対応
  - **[各種UIカスタマイズ](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki/%E3%81%93%E3%81%93%E3%81%8B%E3%82%89%E5%A7%8B%E3%82%81%E3%82%88%E3%81%86%EF%BC%81#ui%E3%82%AB%E3%82%B9%E3%82%BF%E3%83%9E%E3%82%A4%E3%82%BA%E4%BB%95%E6%A7%98%E6%9B%B8)**
    - アニメーション付きスコア表示
    - 異なる要素でアニメーション速度を調整する
    - 交換可能なAPコンボ
    - 交換可能な判定タイプ（PERFECT/GREAT/GOODなど）
    - 交換可能なライフ値
    - 達成率表示
  - **追加サーバーの対応：**
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

## 動画例（画像をクリックして視聴）
  - TootieJin

<a href="https://youtu.be/uXx1OZDQZOI"><img width=32% height=32% alt="【Project Sekai x Honkai: Star Rail】Nameless Faces - HOYO-MiX feat. Lilas Ikuta (Fanmade)" src="https://img.youtube.com/vi/uXx1OZDQZOI/maxresdefault.jpg" /></a> <a href="https://youtu.be/BHVNuwxA1ek"><img width=32% height=32% alt="【Project Sekai Fanmade? (v3→v1)】Hello, SEKAI - DECO*27【ETERNAL Lv32】" src="https://img.youtube.com/vi/BHVNuwxA1ek/maxresdefault.jpg" /></a> <a href="https://youtu.be/O2GxrLWbwlY"><img width=32% height=32% alt="【Project Sekai Fanmade】Confessions of a Rotten Girl - SAWTOWNE (APPEND Lv36)" src="https://img.youtube.com/vi/O2GxrLWbwlY/maxresdefault.jpg" /></a>
<a href="https://youtu.be/yyoTAJx-6ok"><img width=32% height=32% alt="【Project Sekai Fanmade (?)】Anonymous M - PinocchioP (CO-OP CHART)" src="https://img.youtube.com/vi/yyoTAJx-6ok/maxresdefault.jpg" /></a>

  - Hydrogen

<a href="https://youtu.be/aXAeV8nggPs"><img width=32% height=32% alt="【創作譜面】思いついてしまったんだ...... 曲: オーバーライド - 作: 吉田夜世" src="https://img.youtube.com/vi/aXAeV8nggPs/maxresdefault.jpg" /></a> <a href="https://youtu.be/Jqf5eb8neOA"><img width=32% height=32% alt="【創作譜面】易しいギミック 曲: OCO LIMBO - 作: MYUKKE." src="https://img.youtube.com/vi/Jqf5eb8neOA/maxresdefault.jpg" /></a> <a href="https://youtu.be/zOqBCZKK5IY"><img width=32% height=32% alt="【創作譜面】久々の公式レギュ 曲: カルチャー - 作: キタニタツヤ" src="https://img.youtube.com/vi/zOqBCZKK5IY/maxresdefault.jpg" /></a>

  - まてぃ

<a href="https://youtu.be/D0VD5fdQq3M"><img width=32% height=32% alt="マシュマロ / Marshmallow Master Lv33 【プロセカ創作譜面】" src="https://img.youtube.com/vi/D0VD5fdQq3M/maxresdefault.jpg" /></a> <a href="https://youtu.be/yvdWJUD1D9A"><img width=32% height=32% alt="アベリア APPEND Lv35 【プロセカ創作譜面】" src="https://img.youtube.com/vi/yvdWJUD1D9A/maxresdefault.jpg" /></a>

  - ReiyuN

<a href="https://youtu.be/R4jHeilZ3uk"><img width=32% height=32% alt="【 プロセカ 創作譜面 】wiege | Alien Stage x Project SEKAI【 APPEND 27 】" src="https://img.youtube.com/vi/R4jHeilZ3uk/maxresdefault.jpg" /></a> <a href="https://youtu.be/PTuaPc2vSyo"><img width=32% height=32% alt="【 Project SEKAI Fanmade 】Static【 MASTER 30 】" src="https://img.youtube.com/vi/PTuaPc2vSyo/maxresdefault.jpg" /></a>

  - Runakkushia

<a href="https://youtu.be/X4HoJvo6xRE"><img width=32% height=32% alt="【Project Sekai Fanmade】Don’t Believe in T feat. Hatsune Miku & Kasane Teto - PinocchioP" src="https://img.youtube.com/vi/X4HoJvo6xRE/maxresdefault.jpg" /></a>

## 利用規約

1. **(必須)** 動画の説明文に、こちらのテキストをコピーしてください：

**EN**
```
PJSekai-style video creation tool:
- Forked ver. by TootieJin (https://tootiejin.com) & ぴぃまん (https://pim4n-net.com)
- Developed by 名無し｡ (https://sevenc7c.com) 
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

**JP**
```
プロセカ風動画作成補助ツール：
- フォーク：TootieJin氏 (https://tootiejin.com) & ぴぃまん氏 (https://pim4n-net.com)
- 作成：名無し｡氏 (https://sevenc7c.com)
   https://github.com/TootieJin/pjsekai-overlay-APPEND
```

2. 決して**悪意のある使用**をしないでください。（SNS上でデマを流すために使う等）
3. このツールを使ったことによるトラブルや不利益などが発生しても、作者は**一切の責任を負いません。**

## 必須事項
| AviUtl                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | **AviUtl ExEdit2 (推奨)**                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [拡張編集プラグイン](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip)<br>- **インストーラー：** [AviUtl インストーラースクリプト](https://github.com/menndouyukkuri/aviutl-installer-script)<br><br>*すべてのプラグインを `aviutl\Plugins` フォルダにインポートする* | [AviUtl ExEdit2 **(インストーラ版)**](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [MP4Exporter](https://apps.esugo.net/aviutl2-mp4exporter/)<br><br>*`AviUtl2\lwinput.aui2`と`MP4Exporter.auo2`プラグインを`C:\ProgramData\aviutl2\Plugin`フォルダにインポートする* |
- **フォント：** ロダンNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- システムロケールは **必ず「日本語（日本）」** に設定してください [(変更方法はこちら)](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki/%E3%81%93%E3%81%93%E3%81%8B%E3%82%89%E5%A7%8B%E3%82%81%E3%82%88%E3%81%86%EF%BC%81#%E3%82%B7%E3%82%B9%E3%83%86%E3%83%A0%E3%83%AD%E3%82%B1%E3%83%BC%E3%83%AB%E3%81%AE%E5%A4%89%E6%9B%B4)
- **日本語言語パック**をインストールする必要があります [(インストール手順はこちら)](https://support.microsoft.com/ja-jp/windows/windows-%E7%94%A8%E8%A8%80%E8%AA%9E%E3%83%91%E3%83%83%E3%82%AF-a5094319-a92d-18de-5b53-1cfc697cfca8)
- AviUtlの基本的な知識

## 利用方法

**ガイドは[Wikiページ](https://github.com/TootieJin/pjsekai-overlay-APPEND/wiki)に移動しました。**

> [!NOTE]
> 情報を更新する場合がありますので、Wikiを頻繁に確認してください。

## 特別なお礼
- **[@sevenc-nanashi](https://github.com/sevenc-nanashi)氏による[pjsekai-overlay](https://github.com/sevenc-nanashi/pjsekai-overlay)と[pjsekai-background-gen-rust](https://github.com/sevenc-nanashi/pjsekai-background-gen-rust)ツールの開発に感謝します！**
- [@Piliman22](https://github.com/Piliman22)氏、pjsekai-overlayにおける[ScoreSync](https://github.com/Piliman22/ScoreSync)の実現に貢献いただき感謝申し上げます！
- [@Reiyunkun](https://github.com/Reiyunkun)氏、[@gaven1880](https://github.com/gaven1880)氏、[@YumYummity](https://github.com/YumYummity)氏による[追加のプロセカアセット](./extra-assets)の提供に感謝します！
- [@Khronophobia](https://github.com/Khronophobia)氏にはBlenderでのカスタムレーンアセットを提供いただきました！
- [@nagotown](https://github.com/nagotown)氏による[aviutl2-EN](https://github.com/nagotown/aviutl2-EN)翻訳スクリプト！
- [@MattMayuga](https://github.com/MattMayuga)氏による[カスタマイズされた判定フォント](https://github.com/Tiny-Foxes/JudgeFonts-by-MattMayuga)！
- そして私のツールを使ってくださった皆様、本当にありがとうございます。
