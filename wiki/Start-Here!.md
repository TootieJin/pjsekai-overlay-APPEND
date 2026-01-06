# Requirements

| AviUtl                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | **AviUtl ExEdit2 (Recommended)**                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [AviUtl](http://spring-fragrance.mints.ne.jp/aviutl/) + [Advanced Editing plug-in](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [easymp4](https://aoytsk.blog.jp/aviutl/easymp4.zip) ([JP Installation Guide](https://aviutl.info/dl-innsuto-ru/))<br>- **JP installer:** [AviUtl JP Installer Script](https://github.com/menndouyukkuri/aviutl-installer-script)<br>- **EN installer:** [AviUtl EN Extra Pack](https://www.videohelp.com/download/AviUtl_setup_1.14.exe)<br><br>*Import all plugins to `aviutl\Plugins` folder* | [AviUtl ExEdit2 **(via installer)**](http://spring-fragrance.mints.ne.jp/aviutl/) + [L-SMASH Works](https://github.com/Mr-Ojii/L-SMASH-Works-Auto-Builds/releases/latest) + [MP4Exporter](https://apps.esugo.net/aviutl2-mp4exporter/)<br>- English language can be selected in `設定 > 言語の設定 > English` (`Settings > Language > English`)<br><br>*Import `AviUtl2\lwinput.aui2` & `MP4Exporter.auo2` plugin to `C:\ProgramData\aviutl2\Plugin` folder* |
- **[PowerShell](https://github.com/PowerShell/PowerShell)**
- **Fonts:** RodinNTLG [DB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-db/) + [EB](https://en.fontworks.co.jp/fontsearch/rodinntlgpro-eb/)
- Your system locale **must be `Japanese (Japan)`** [(Go here for how to change it)](#changing-system-locale)
- **Japanese Language Pack** must be installed [(Go here for installation instructions)](https://support.microsoft.com/en-us/windows/language-packs-for-windows-a5094319-a92d-18de-5b53-1cfc697cfca8)
- Basic knowledge of AviUtl

# Video Guide

You'll need **a clean recording that has BLACK background & ALL UI hidden (including lanes)** like in the screenshot below:

<img width=60% height="auto" alt="" src="https://github.com/user-attachments/assets/46578e52-a3e6-4248-9e2c-f61486cc6af6" />

Follow the steps below:
1. Go to [Sonolus](https://sonolus.com/) to find your chart.
2. Screen record the video with **BLACK Background**, **`Stage` OFF**, **`Hide UI` ON or ALL** with **`Note Margin` >= 0.08**

> [!NOTE]
> To change the background, go to `Configuration > Background > Select`
> 
> <img width=60% height="auto" alt="" src="https://github.com/user-attachments/assets/a138e3b7-50cd-49e0-8977-c1df60ba9a3e" />

3. Transfer the video file to your computer.
4. **(AviUtl ONLY)** Download the [FFmpeg](https://www.ffmpeg.org/) encoder and use it (`ffmpeg -i source.mp4 output.mp4`)
   - **This step is required because of a funny bug only seen in AviUtl**: If the frame rate of the source video does not match the project, synchronization is delayed.
   - Alternatively, you can use your video editor to do this step.
5. Once done, refer to the usage guide below.

# Usage Guide
0. [Install AviUtl & other plugins](#requirements)
1. Download the latest version of pjsekai-overlay-APPEND [here](https://github.com/TootieJin/pjsekai-overlay-APPEND/releases/latest/).
   - **(optional)** You can download `extra-assets.zip` for use of non-essential assets.
2. Unzip the file in a directory that **doesn't contain non-ASCII characters nor require administrative permissions** (see [#5](https://github.com/TootieJin/pjsekai-overlay-APPEND/issues/5))
3. Follow the steps based on which AviUtl you want to use:

| AviUtl                                                                                                                                                                                                                                                                                                                           	| **AviUtl ExEdit2**                                                                      	|
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------	|
| 1. **While `aviutl.exe` is open**, open `pjsekai-overlay-APPEND.exe`<br>2. **Restart AviUtl to apply changes.** 	| 1. Open `pjsekai-overlay-APPEND.exe`<br>2. Press `2` to choose the AviUtl ExEdit2 instance 	|

4. In your `pjsekai-overlay-APPEND.exe` command prompt, input the chart ID including the prefix.
   - **Official servers:**
      - `sekai-rush-`: Proseka Rush ([sekairush.shop](https://sekairush.shop))
   - **Custom servers:**
      - `chcy-`: Chart Cyanvas
      - `0`: Archive ([cc.sevenc7c.com](https://cc.sevenc7c.com))
      - `1`: Offshoot server ([chart-cyanvas.com](https://chart-cyanvas.com))
      - `Others (URL domain)`: Different Cyanvas instance
      - `ptlv-`: Potato Leaves ([ptlv.sevenc7c.com](https://ptlv.sevenc7c.com))
      - <del>`utsk-`: Untitled Sekai ([us.pim4n-net.com](https://us.pim4n-net.com))</del>
      - `UnCh-`: UntitledCharts ([untitledcharts.com](https://untitledcharts.com))
      - `coconut-next-sekai-`: Next SEKAI ([coconut.sonolus.com/next-sekai](https://coconut.sonolus.com/next-sekai))
      - `lalo-`: [laoloser](https://www.youtube.com/@laoloserr)'s server ([sonolus.laoloser.com](https://laoloser.com/sonolus))
      - `sync-`: Local Server ([ScoreSync](https://github.com/Piliman22/ScoreSync))

5. Create a new project in AviUtl with **Image size 1920x1080 (or 1440x1080) & Frame rate 60**

| AviUtl                                                                                                                                                                                                                                                                                                                           	| **AviUtl ExEdit2**                                                                      	|
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------	|
| Right click anywhere in the `Advanced Editing` window > `Create a New Project`<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/7b757370-e575-4ba2-b1f9-752da586d938" /> 	| Go to `File > New Project`<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/253c9b76-6194-4fc4-8649-f711f230fcbc" /> 	|

6. Import specified exo/alias(.object) file by navigating to your `pjsekai-overlay-APPEND\dist\[Chart ID]` directory:

|      	| AviUtl                                                          	                                                                                            | **AviUtl ExEdit2**                                                	|
|------	|--------------------------------------------------------------------------------------------------------------------------------------------------------------	|-------------------------------------------------------------------- |
| 16:9 	| - **EN ver:**<br>`main_en_16-9_1920x1080.exo`<br>`v1-skin_en_16-9_1920x1080.exo`<br>- **JP ver:**<br>`main_jp_16-9_1920x1080.exo`<br>`v1-skin_jp_16-9_1920x1080.exo` 	| `main2_16-9_1920x1080.object`<br>`v1-skin2_16-9_1920x1080.object` 	|
| 4:3  	| - **EN ver:**<br>`main_en_4-3_1440x1080.exo`<br>`v1-skin_en_4-3_1440x1080.exo`<br>- **JP ver:**<br>`main_jp_4-3_1440x1080.exo`<br>`v1-skin_jp_4-3_1440x1080.exo`     	| `main2_4-3_1440x1080.object`<br>`v1-skin2_4-3_1440x1080.object`   	|

7. **Select a video file in the timeline** provided by the exo/alias(.object)

| AviUtl                                                                                                                                                                                                                                                                                                                           	| **AviUtl ExEdit2**                                                                      	|
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|--------------------------------------------------------------------------------------------	|
| 1. Double click `Video file[Standard drawing]`<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/37dda2f7-364d-4945-889d-51c2fc3b74de" /><br><br>2. Click `Reference file`, then select a video file<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/7d88eb76-f682-44c4-a7ee-bcf36ddf2075" /> 	| 1. Click `Video file`<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/6d906632-0c1c-444c-8aa6-a0447137dfb8" /><br><br>2. Click `File`, then select a video file<br><br><img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/b790401c-8784-4dc8-9f0b-363164970013" /> 	|

8. Adjust the video & audio positioning as necessary. **Make sure the video frame matches the slot particle with the lane**

<img width=50% height="auto" alt="image" src="https://github.com/user-attachments/assets/78b91894-57c7-424a-af2d-3ee6f2ff8fb8" />

9. Go to `Root@pjsekai-overlay` and **adjust the `Offset`**

<img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/433a1e1f-a24b-4518-8526-e73192eaf38a" />

10. Once finished, export your video as mp4
   - **AviUtl:** `File > Export with Plugin > かんたんMP4出力` or `File > Export with Plugin > Adv. x264/x265 Export(GUI) Ex`
   - **AviUtl ExEdit2:** `File > Export > MP4 Exporter (by えすご/Esugo)`

# UI Customization Specifications
> [!TIP]
> You can change default values in `default.ini`, **reducing the need to change values in AviUtl.**
> 
> <img width=70% height=70% alt="image" src="https://github.com/user-attachments/assets/bd2151ad-2045-46ef-96c2-cbbb426a68d9" />

### Root@pjsekai-overlay-en
<img width="383" height="79" alt="image" src="https://github.com/user-attachments/assets/76da6c77-f6a7-4480-b279-b5d53f3e583f" />

| **Name**      	| Description                                                                                                	| Default 	|        Range       	|
|---------------	|------------------------------------------------------------------------------------------------------------	|:-------:	|:------------------:	|
| **Offset**<br>`offset`    	| Number of frames to shift events<br>- Increase to shift timing later<br>- Decrease to shift timing earlier 	|  216.00  	| -99999.99 ~ 99999.99 	|
| **Cache**<br>`cache`     	| When cache is set to 0, any change in the `data.ped` is applied immediately<br>(`0` - OFF, `1` - ON)                               	|    1    	|        0 or 1       	|
| **Text Language**<br>`text_lang` 	| Set text language configuration for the error & detailed stats text<br>(`0` - 日本語, `1` - English)                 	|    0    	|        0 or 1       	|
| **Watermark**<br>`watermark` 	| Enable watermark text at the bottom-left corner                                                            	|   1 (ON)  	|          0 or 1         	|
| **Detailed Stats**<br>`detail_stat` 	| <img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/9e916e51-d888-497e-b649-b6e9b4a3b760" /><br>View detailed statistics (useful for measuring & analysing). _This text will not be shown when exporting_                 	|    0 (OFF)    	|        0 or 1       	|

### Life@pjsekai-overlay-en
<img width="125" height="125" alt="LifeUP" src="https://github.com/user-attachments/assets/6f7a7db8-50bb-43cf-9463-5f46325c862e" /> <img width=30% height=30% alt="life-v1" src="https://github.com/user-attachments/assets/b467f932-fa04-4afa-bbf2-9aed811b0855" /> <img width=30% height=30% alt="life" src="https://github.com/user-attachments/assets/7aab3534-66cf-4dad-936e-3d423ecce615" />

| **Name** 	| Description                                                                         	| Default 	|   Range  	|
|----------	|-------------------------------------------------------------------------------------	|:-------:	|:--------:	|
| **LIFE**<br>`life` 	| LIFE value (self-explanatory)<br>- When value changes, the LIFE bar changes as well 	|   1000  	| 0 ~ 9999 (Integer) 	|
| **Skill Effect**<br>`life_skill` 	| Toggle skill glow effect 	|  0 (OFF)  	|    0 or 1   	|
| **Overflow LIFE Bar**<br>`overflow` 	| <img width=70% height=70% alt="life_overflow-v1" src="https://github.com/user-attachments/assets/2fd69cf1-a767-47e7-adc3-14ecd2d56ce6" /><br><img width=70% height=70% alt="life_overflow" src="https://github.com/user-attachments/assets/bbe6ea47-b8e8-498f-8e30-187ca8971c64" /> 	|  0 (OFF)  	|    0 or 1   	|
| **Leading Zero**<br>`lead_zero` 	| Append the "0" digits for LIFE value below 1000 	|  0 (OFF)  	|    0 or 1   	|

### Score@pjsekai-overlay-en
<img width="125" height="125" alt="ScoreUP" src="https://github.com/user-attachments/assets/a5a8b0f0-035c-4951-8ae3-d2038945d86c" /> <img width=50% height=50% alt="bg" src="https://github.com/user-attachments/assets/3db93b3e-2280-46e1-a08f-00e50a5e5e8c" />

| **Name**             	| Description                                    	| Default 	|  Range 	|
|----------------------	|------------------------------------------------	|:-------:	|:------:	|
| **Min Digit**<br>`min_digit`        	| Render the minimum amount of digits in score   	|    8    	| 1 ~ 99 (Integer)	|
| **Skill Effect**<br>`score_skill`        	| Toggle skill glow effect<br>`0`: OFF<br>`1`: AUTO - For each skill event, show skill glow effect for 5 seconds<br>`2`: ON   	|    1  	| 0 ~ 2 (Integer)	|
| **Animation Speed**<br>`score_speed`   | Adjust animation speed                       	   |   1.00     |    >= 0 	   |
| **Animated Scoring**<br>`anim_score` 	| Increase incrementally rather than all at once 	|  0 (OFF)  	|    0 or 1   	|
| **WDS animation**<br>`wds_anim` 	| Toggle World Dai Star's added score animation 	|  0 (OFF)  	|    0 or 1   	|

> [!NOTE]
> **About skill events: https://sekai-guide.tootiejin.com/charting-guide/adding-skill-and-fever-events**

> [!TIP]
> **You can add skill events to use the AUTO Skill Effect.** Go to your `data.ped` file and for each skill event you want to add, add the line `s|[timeframe (seconds)]`. You can add this as many as you want.
> - **Example:**
> 
> <img width=50% height=50% alt="image" src="https://github.com/user-attachments/assets/44bb1ec6-5d42-467c-af1b-53b583d8ad44" />


### Combo@pjsekai-overlay-en
<img width="148" height="49" alt="pt" src="https://github.com/user-attachments/assets/9db50558-cf81-4ed8-a2bd-1d4bbd22e156" /> <img width="145" height="45" alt="nt" src="https://github.com/user-attachments/assets/3ca0f65e-8ce6-40c9-8ff4-53a8ee9d2f81" />

| **Name**                          	| Description                            	  | Default | Range 	|
|-----------------------------------	|----------------------------------------	  |:-------:|:-----:	|
| **AP Combo**<br>`ap`                      	| Toggle AP Combo status                 	  |    1   	| 0 or 1 	|
| **Combo Tag**<br>`tag`                     	| Toggle rendering combo tag             	  |    1   	| 0 or 1 	|
| **Show the last X digit(s)**<br>`last_digit`        	| `4`: 12345 -> /2345<br>`2`: 12345 -> ///45 |    4   	|   >= 0 (Integer)   	|
| **Animation Speed**<br>`combo_speed`   | Adjust animation speed                  |   1.00     |    >= 0 	   |
| **n00 Combo Burst**<br>`combo_burst`        	| Toggle World Dai Star's combo burst effect |    0 (OFF)   	|   0 or 1   	|
| **Achievement Rate**<br>`achievement_rate`        	| <img width="auto" height="auto" alt="image" src="https://github.com/user-attachments/assets/1e1cc392-fc5d-4c7d-b954-a6eba95d8839" /><br>Toggle Achievement Rate displayed on top of the combo counter **(Set to `0` or leave it blank to turn it OFF)** <br><br>**Extra features:**<br>`100.0000+`: Increase incrementally (like Animated Scoring) up to 100.0000%<br>`100.0000-`: Start at 100.0000%, then decrease for non-perfect judgements |    100.0000   	|   Valid floating number   	|

### Judgement@pjsekai-overlay-en
<img width="125" height="125" alt="SkillUP" src="https://github.com/user-attachments/assets/e29f426d-71ae-4de5-912a-a5c7375f538d" />

| **Name**       	| Description                                                                   	| Default 	| Range 	|
|----------------	|-------------------------------------------------------------------------------	|:-------:	|:-----:	|
| **Judge Type**<br>`judge` 	| `1`: <img width=25% height=25% alt="perfect" src="https://github.com/user-attachments/assets/28950e9e-0dac-49d9-81d3-70bdaa2d6f0c" /><br>`2`: <img width=25% height=25% alt="great" src="https://github.com/user-attachments/assets/ccf333a7-795d-43ad-8002-a9d2220e18a6" /><br>`3`: <img width=25% height=25% alt="good" src="https://github.com/user-attachments/assets/9d0a26bb-c8e7-47d0-9a3d-717b4ad0e0fa" /><br>`4`: <img width=25% height=25% alt="bad" src="https://github.com/user-attachments/assets/5b757195-8bd4-4beb-9f77-808000f1d865" /><br>`5`: <img width=25% height=25% alt="miss" src="https://github.com/user-attachments/assets/734ead15-491b-4bdb-9017-f2b30ab32223" /><br>`6`: <img width=25% height=25% alt="auto" src="https://github.com/user-attachments/assets/b9d674cf-1b69-478e-b2be-53691109b12d" /><br>`7 ~ 10`: Custom judgements 	|    1    	| 1 ~ 10 (Integer) 	|
| **Animation Speed**<br>`judge_speed`                 	| Adjust animation speed               	  |   1.00   |   >= 0 	|

# Common Problems
### Changing system locale
Change your system locale to `Japanese (Japan)`: `Settings > Time & Language > Language > Administrative language settings > "Administrative" tab > Change system locale...`

<img width=40% height=40% alt="image" src="https://github.com/user-attachments/assets/0b5b680a-56cf-4cdb-8bb9-f10a50d7c07a" />

### Installing a language pack
- [Go here for installation instructions.](https://support.microsoft.com/en-us/windows/language-packs-for-windows-a5094319-a92d-18de-5b53-1cfc697cfca8)

### The animation slowing down when importing
Go to `File > New Project` and set frame rate to 60. The change will be reflected in future imports.

<img width="418" height="178" alt="image" src="https://github.com/user-attachments/assets/6bff30c3-fb67-4169-931f-b3df7c7d65f6" />

### The UI doesn't show anything / PED file unable to be opened
You likely haven't selected a **`data.ped` file.** 
- **AviUtl:** Locate your `data.ped` in the `pjsekai-overlay-APPEND\dist\[Chart ID]` directory & copy the path to the `Root@pjsekai-overlay-en > Setting > PED File Path`. **(Make sure to replace `\ (¥)` with either `\\ (¥¥)` or `/`).**

<img width=50% height="auto" alt="image" src="https://github.com/user-attachments/assets/d34c198d-c819-4eea-a545-27aca04eba9c" />

- **AviUtl ExEdit2:** Go to `Root@pjsekai-overlay-2`, select `PED File`, then select `data.ped` in the `pjsekai-overlay-APPEND\dist\[Chart ID]` directory

<img width="271" height="94" alt="image" src="https://github.com/user-attachments/assets/a08f3a04-4ac1-41d1-a850-7d08250366d4" />

### The intro doesn't load all objects
AviUtl for whatever reason stops loading when importing for the first time. If that is the case, **you may have to import it again.**

### The fonts look off/incorrect

Project Sekai uses the **RodinNTLG DB + EB** font (described in [Requirements](#requirements)). These are paid fonts so we're not allowed to show you how to install it.

<img width=50% height=50% alt="image" src="https://github.com/user-attachments/assets/9c125788-ab05-4fc0-8d2a-6f3e5d9c892a" />

### Still have questions or encountering problems?
**[Make a discussion thread.](https://github.com/TootieJin/pjsekai-overlay-APPEND/discussions)**