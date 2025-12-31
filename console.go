package main

import (
	"fmt"
	"strings"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/pjsekaioverlay"
	"github.com/lithammer/dedent"
)

func Title() {
	fmt.Printf(
		strings.TrimSpace(dedent.Dedent(`
    %s-- pjsekai-overlay-%sAP%sPE%sND %s-----------------------------------------------------------%s
    %sフォークプロセカ風動画作成補助ツール / Forked PJSekai-style video creation tool%s
        Version: %s%s%s

        Developed by %s名無し｡(@sevenc-nanashi)%s
            https://github.com/sevenc-nanashi/pjsekai-overlay
        Forked by %sTootieJin & ぴぃまん(@Piliman22)%s
         -> https://github.com/TootieJin/pjsekai-overlay-APPEND %s(使用中/In use)%s
	
    %s[CAUTION] This tool is primary only for people with technical know-how and basic knowledge of AviUtl.%s 
    %sIf you have any questions/problems, please make a discussion thread. Refer to the wiki for how to set it up.%s

    %s[注意] このツールは主に、技術的な知識とAviUtlの基本的な理解がある方のみを対象としています。%s 
    %s質問や問題がある場合は、議論スレッドを作成してください。設定方法についてはWikiを参照してください。%s
    %s-------------------------------------------------------------------------------------%s
    `))+"\n",
		RgbColorEscape(0x00afc7), RgbColorEscape(0xab93ff), RgbColorEscape(0xd388ed), RgbColorEscape(0xff8bf4), RgbColorEscape(0x00afc7), ResetEscape(),
		RgbColorEscape(0x00afc7), ResetEscape(),
		RgbColorEscape(0x0f6ea3), pjsekaioverlay.Version, ResetEscape(),
		RgbColorEscape(0x48b0d5), ResetEscape(),
		RgbColorEscape(0x48b0d5), ResetEscape(),
		RgbColorEscape(0xadff2f), ResetEscape(),
		RgbColorEscape(0xff0000), ResetEscape(),
		RgbColorEscape(0xff0000), ResetEscape(),
		RgbColorEscape(0xff0000), ResetEscape(),
		RgbColorEscape(0xff0000), ResetEscape(),
		RgbColorEscape(0xff5a91), ResetEscape(),
	)

}

func RgbColorEscape(rgb int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", (rgb>>16)&0xff, (rgb>>8)&0xff, rgb&0xff)
}

func AnsiColorEscape(color int) string {
	return fmt.Sprintf("\033[38;5;%dm", color)
}

func ResetEscape() string {
	return "\033[0m"
}
