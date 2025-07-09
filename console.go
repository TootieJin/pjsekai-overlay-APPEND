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
    %s== pjsekai-overlay-%sAP%sPE%sND %s-----------------------------------------------------------%s
    %sフォークプロセカ風動画作成補助ツール / Forked PJSekai-style video creation tool%s
        Version: %s%s%s

        Developed by %s名無し｡(@sevenc-nanashi)%s
            https://github.com/sevenc-nanashi/pjsekai-overlay
        Forked by %sTootieJin%s
         -> https://github.com/TootieJin/pjsekai-overlay-APPEND %s(使用中/In use)%s
            https://github.com/TootieJin/pjsekai-overlay-MULTI
	
    %s[CAUTION] This tool is primary only for people with technical know-how and basic knowledge of AviUtl.%s 
    %sDO NOT open issues or request help about this. Refer to the README for how to set it up.%s
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
