package pjsekaioverlay

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func encodeString(str string) string {
	bytes := utf16.Encode([]rune(str))
	encoded := make([]string, 1024)
	if len(str) > 1024 {
		panic("too long string")
	}
	for i := range encoded {
		var hex string
		if i >= len(bytes) {
			hex = fmt.Sprintf("%04x", 0)
		} else {
			hex = fmt.Sprintf("%02x%02x", bytes[i]&0xff, bytes[i]>>8)
		}

		encoded[i] = hex
	}

	return strings.Join(encoded, "")
}

//go:embed main_jp_16-9_1920x1080.exo
var rawBaseExoJP []byte

//go:embed main_jp_4-3_1440x1080.exo
var rawBaseExoJP43 []byte

//go:embed main_en_16-9_1920x1080.exo
var rawBaseExoEN []byte

//go:embed main_en_4-3_1440x1080.exo
var rawBaseExoEN43 []byte

//go:embed v1-skin_jp_16-9_1920x1080.exo
var rawBaseExoJPv1 []byte

//go:embed v1-skin_jp_4-3_1440x1080.exo
var rawBaseExoJP43v1 []byte

//go:embed v1-skin_en_16-9_1920x1080.exo
var rawBaseExoENv1 []byte

//go:embed v1-skin_en_4-3_1440x1080.exo
var rawBaseExoEN43v1 []byte

//go:embed main2_16-9_1920x1080.object
var rawBaseAlias []byte

//go:embed main2_4-3_1440x1080.object
var rawBaseAlias43 []byte

//go:embed v1-skin2_16-9_1920x1080.object
var rawBaseAliasv1 []byte

//go:embed v1-skin2_4-3_1440x1080.object
var rawBaseAlias43v1 []byte

func WriteExoFiles(assets string, destDir string, title string, description []string, descriptionv1 []string, difficulty string, extra string, exFile string, exFileOpacity string, mappingFile []string) error {
	baseExoJP := string(rawBaseExoJP)
	baseExoJP43 := string(rawBaseExoJP43)
	baseExoEN := string(rawBaseExoEN)
	baseExoEN43 := string(rawBaseExoEN43)
	baseExoJPv1 := string(rawBaseExoJPv1)
	baseExoJP43v1 := string(rawBaseExoJP43v1)
	baseExoENv1 := string(rawBaseExoENv1)
	baseExoEN43v1 := string(rawBaseExoEN43v1)

	mapping := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", encodeString(difficulty),
		"{text:extra}", encodeString(extra),
		"{text:title}", encodeString(title),
		"{text:description-1}", encodeString(description[0]),
		"{text:description-2}", encodeString(description[1]),
		"{image:tournament}", exFile,
		"{opacity}", exFileOpacity,
		"{difficulty}", strings.ToLower(difficulty),
		// Root
		"{offset}", mappingFile[0], // track0
		"{cache}", mappingFile[1], // track1
		"{text_lang}", mappingFile[2], // dialog: text_lang=
		"{watermark}", mappingFile[3], // check0
		"{detail_stat}", mappingFile[4], // dialog: detail_stat=
		// Life
		"{life}", mappingFile[5], // track1
		"{life_skill}", mappingFile[6], // track2
		"{overflow}", mappingFile[7], // dialog: overflow=
		"{lead_zero}", mappingFile[8], // dialog: lead_zero=
		// Score
		"{min_digit}", mappingFile[9], // track1
		"{score_skill}", mappingFile[10], // track2
		"{score_speed}", mappingFile[11], // dialog: speed=
		"{anim_score}", mappingFile[12], // check0
		"{wds_anim}", mappingFile[13], // dialog: wds_anim=
		// Combo
		"{ap}", mappingFile[14], // track1
		"{tag}", mappingFile[15], // track2
		"{last_digit}", mappingFile[16], // dialog: digits=
		"{combo_speed}", mappingFile[17], // dialog: speed=
		"{combo_burst}", mappingFile[18], // dialog: combo_burst=
		"{achievement_rate}", mappingFile[19], // dialog: achievement_rate=
		// Judgement
		"{judge}", mappingFile[20], // track0
		"{judge_speed}", mappingFile[21], // dialog: speed=
	}

	mappingv1 := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", encodeString(difficulty),
		"{text:extra}", encodeString(extra),
		"{text:title}", encodeString(title),
		"{text:description-1}", encodeString(descriptionv1[0]),
		"{text:description-2}", encodeString(descriptionv1[1]),
		"{image:tournament}", exFile,
		"{opacity}", exFileOpacity,
		"{difficulty}", strings.ToLower(difficulty),
		// Root
		"{offset}", mappingFile[0], // track0
		"{cache}", mappingFile[1], // track1
		"{text_lang}", mappingFile[2], // dialog: text_lang=
		"{watermark}", mappingFile[3], // check0
		"{detail_stat}", mappingFile[4], // dialog: detail_stat=
		// Life
		"{life}", mappingFile[5], // track1
		"{life_skill}", mappingFile[6], // track2
		"{overflow}", mappingFile[7], // dialog: overflow=
		"{lead_zero}", mappingFile[8], // dialog: lead_zero=
		// Score
		"{min_digit}", mappingFile[9], // track1
		"{score_skill}", mappingFile[10], // track2
		"{score_speed}", mappingFile[11], // dialog: speed=
		"{anim_score}", mappingFile[12], // check0
		"{wds_anim}", mappingFile[13], // dialog: wds_anim=
		// Combo
		"{ap}", mappingFile[14], // track1
		"{tag}", mappingFile[15], // track2
		"{last_digit}", mappingFile[16], // dialog: digits=
		"{combo_speed}", mappingFile[17], // dialog: speed=
		"{combo_burst}", mappingFile[18], // dialog: combo_burst=
		"{achievement_rate}", mappingFile[19], // dialog: achievement_rate=
		// Judgement
		"{judge}", mappingFile[20], // track0
		"{judge_speed}", mappingFile[21], // dialog: speed=
	}
	for i := range mapping {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(baseExoJP, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(baseExoJP43, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(baseExoEN, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(baseExoEN43, mapping[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate exo file) [Missing: %s]", mapping[i-1]))
		}
		baseExoJP = strings.ReplaceAll(baseExoJP, mapping[i-1], mapping[i])
		baseExoJP43 = strings.ReplaceAll(baseExoJP43, mapping[i-1], mapping[i])
		baseExoEN = strings.ReplaceAll(baseExoEN, mapping[i-1], mapping[i])
		baseExoEN43 = strings.ReplaceAll(baseExoEN43, mapping[i-1], mapping[i])
	}
	for i := range mappingv1 {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(baseExoJPv1, mappingv1[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate v1 exo file) [Missing: %s]", mappingv1[i-1]))
		}
		if !strings.Contains(baseExoJP43v1, mappingv1[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate v1 exo file) [Missing: %s]", mappingv1[i-1]))
		}
		if !strings.Contains(baseExoENv1, mappingv1[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate v1 exo file) [Missing: %s]", mappingv1[i-1]))
		}
		if !strings.Contains(baseExoEN43v1, mappingv1[i-1]) {
			panic(fmt.Sprintf("exoファイルの生成に失敗しました (Failed to generate v1 exo file) [Missing: %s]", mappingv1[i-1]))
		}
		baseExoJPv1 = strings.ReplaceAll(baseExoJPv1, mappingv1[i-1], mappingv1[i])
		baseExoJP43v1 = strings.ReplaceAll(baseExoJP43v1, mappingv1[i-1], mappingv1[i])
		baseExoENv1 = strings.ReplaceAll(baseExoENv1, mappingv1[i-1], mappingv1[i])
		baseExoEN43v1 = strings.ReplaceAll(baseExoEN43v1, mappingv1[i-1], mappingv1[i])
	}

	encodedExoJP, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoJP), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoJP43, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoJP43), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEN, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoEN), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEN43, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoEN43), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoJPv1, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoJPv1), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoJP43v1, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoJP43v1), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoENv1, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoENv1), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}
	encodedExoEN43v1, err := io.ReadAll(transform.NewReader(
		strings.NewReader(baseExoEN43v1), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return fmt.Errorf("エンコードに失敗しました (Encoding failed) [%w]", err)
	}

	if err := os.WriteFile(filepath.Join(destDir, "main_jp_16-9_1920x1080.exo"),
		encodedExoJP,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_jp_4-3_1440x1080.exo"),
		encodedExoJP43,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_en_16-9_1920x1080.exo"),
		encodedExoEN,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main_en_4-3_1440x1080.exo"),
		encodedExoEN43,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin_jp_16-9_1920x1080.exo"),
		encodedExoJPv1,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin_jp_4-3_1440x1080.exo"),
		encodedExoJP43v1,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin_en_16-9_1920x1080.exo"),
		encodedExoENv1,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin_en_4-3_1440x1080.exo"),
		encodedExoEN43v1,
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	return nil
}

func WriteAliasFiles(assets string, destDir string, title string, description []string, descriptionv1 []string, difficulty string, extra string, exFile string, exFileOpacity string, mappingFile []string) error {
	baseAlias := string(rawBaseAlias)
	baseAlias43 := string(rawBaseAlias43)
	baseAliasv1 := string(rawBaseAliasv1)
	baseAlias43v1 := string(rawBaseAlias43v1)

	var textLang string
	if mappingFile[2] == "1" {
		textLang = "English"
	} else {
		textLang = "日本語"
	}

	mapping := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", difficulty,
		"{text:extra}", extra,
		"{text:title}", title,
		"{text:description-1}", description[0],
		"{text:description-2}", description[1],
		"{image:tournament}", exFile,
		"{opacity}", exFileOpacity,
		"{difficulty}", strings.ToLower(difficulty),
		// Root
		"{offset}", mappingFile[0],
		"{cache}", mappingFile[1],
		"{text_lang}", textLang,
		"{watermark}", mappingFile[3],
		"{detail_stat}", mappingFile[4],
		// Life
		"{life}", mappingFile[5],
		"{life_skill}", mappingFile[6],
		"{overflow}", mappingFile[7],
		"{lead_zero}", mappingFile[8],
		// Score
		"{min_digit}", mappingFile[9],
		"{score_skill}", mappingFile[10],
		"{score_speed}", mappingFile[11],
		"{anim_score}", mappingFile[12],
		"{wds_anim}", mappingFile[13],
		// Combo
		"{ap}", mappingFile[14],
		"{tag}", mappingFile[15],
		"{last_digit}", mappingFile[16],
		"{combo_speed}", mappingFile[17],
		"{combo_burst}", mappingFile[18],
		"{achievement_rate}", mappingFile[19],
		// Judgement
		"{judge}", mappingFile[20],
		"{judge_speed}", mappingFile[21],
	}

	mappingv1 := []string{
		"{assets}", strings.ReplaceAll(assets, "\\", "/"),
		"{dist}", strings.ReplaceAll(destDir, "\\", "/"),
		"{text:difficulty}", difficulty,
		"{text:extra}", extra,
		"{text:title}", title,
		"{text:description-1}", description[0],
		"{text:description-2}", description[1],
		"{image:tournament}", exFile,
		"{opacity}", exFileOpacity,
		"{difficulty}", strings.ToLower(difficulty),
		// Root
		"{offset}", mappingFile[0],
		"{cache}", mappingFile[1],
		"{text_lang}", textLang,
		"{watermark}", mappingFile[3],
		"{detail_stat}", mappingFile[4],
		// Life
		"{life}", mappingFile[5],
		"{life_skill}", mappingFile[6],
		"{overflow}", mappingFile[7],
		"{lead_zero}", mappingFile[8],
		// Score
		"{min_digit}", mappingFile[9],
		"{score_skill}", mappingFile[10],
		"{score_speed}", mappingFile[11],
		"{anim_score}", mappingFile[12],
		"{wds_anim}", mappingFile[13],
		// Combo
		"{ap}", mappingFile[14],
		"{tag}", mappingFile[15],
		"{last_digit}", mappingFile[16],
		"{combo_speed}", mappingFile[17],
		"{combo_burst}", mappingFile[18],
		"{achievement_rate}", mappingFile[19],
		// Judgement
		"{judge}", mappingFile[20],
		"{judge_speed}", mappingFile[21],
	}
	for i := range mapping {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(baseAlias, mapping[i-1]) {
			panic(fmt.Sprintf("aliasファイルの生成に失敗しました (Failed to generate alias file) [Missing: %s]", mapping[i-1]))
		}
		if !strings.Contains(baseAlias43, mapping[i-1]) {
			panic(fmt.Sprintf("aliasファイルの生成に失敗しました (Failed to generate alias file) [Missing: %s]", mapping[i-1]))
		}
		baseAlias = strings.ReplaceAll(baseAlias, mapping[i-1], mapping[i])
		baseAlias43 = strings.ReplaceAll(baseAlias43, mapping[i-1], mapping[i])
	}
	for i := range mappingv1 {
		if i%2 == 0 {
			continue
		}
		if !strings.Contains(baseAliasv1, mappingv1[i-1]) {
			panic(fmt.Sprintf("aliasファイルの生成に失敗しました (Failed to generate alias file) [Missing: %s]", mappingv1[i-1]))
		}
		if !strings.Contains(baseAlias43v1, mappingv1[i-1]) {
			panic(fmt.Sprintf("aliasファイルの生成に失敗しました (Failed to generate alias file) [Missing: %s]", mappingv1[i-1]))
		}
		baseAliasv1 = strings.ReplaceAll(baseAliasv1, mappingv1[i-1], mappingv1[i])
		baseAlias43v1 = strings.ReplaceAll(baseAlias43v1, mappingv1[i-1], mappingv1[i])
	}

	if err := os.WriteFile(filepath.Join(destDir, "main2_16-9_1920x1080.object"),
		[]byte(baseAlias),
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "main2_4-3_1440x1080.object"),
		[]byte(baseAlias43),
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin2_16-9_1920x1080.object"),
		[]byte(baseAliasv1),
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	if err := os.WriteFile(filepath.Join(destDir, "v1-skin2_4-3_1440x1080.object"),
		[]byte(baseAlias43v1),
		0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (Failed to write file) [%w]", err)
	}
	return nil
}
