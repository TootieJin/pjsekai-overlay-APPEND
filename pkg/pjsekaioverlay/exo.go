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

func WriteExoFiles(assets string, destDir string, title string, description []string, descriptionv1 []string, difficulty string, extra string, exFile string, exFileOpacity string, mappingExo []string) error {
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
		"{offset}", mappingExo[0], // track0
		"{cache}", mappingExo[1], // track1
		"{font_type}", mappingExo[2], // track2
		"{watermark}", mappingExo[3], // check0
		// Life
		"{life}", mappingExo[4], // track0
		"{overflow}", mappingExo[5], // check0
		// Score
		"{min_digit}", mappingExo[6], // track0
		"{anim_score}", mappingExo[7], // check0
		"{score_speed}", mappingExo[8], // dialog: speed=
		// Combo
		"{ap}", mappingExo[9], // track0
		"{tag}", mappingExo[10], // track1
		"{last_digit}", mappingExo[11], // dialog: digits=
		"{combo_speed}", mappingExo[12], // dialog: speed=
		"{combo_burst}", mappingExo[13], // dialog: combo_burst=
		// Judgement
		"{judge}", mappingExo[14], // track0
		"{judge_speed}", mappingExo[15], // dialog: speed=
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
		"{offset}", mappingExo[0], // track0
		"{cache}", mappingExo[1], // track1
		"{font_type}", mappingExo[2], // track2
		"{watermark}", mappingExo[3], // check0
		// Score
		"{min_digit}", mappingExo[6], // track0
		"{anim_score}", mappingExo[7], // check0
		"{score_speed}", mappingExo[8], // dialog: speed=
		// Combo
		"{ap}", mappingExo[9], // track0
		"{tag}", mappingExo[10], // track1
		"{last_digit}", mappingExo[11], // dialog: digits=
		"{combo_speed}", mappingExo[12], // dialog: speed=
		"{combo_burst}", mappingExo[13], // dialog: combo_burst=
		// Judgement
		"{judge}", mappingExo[14], // track0
		"{judge_speed}", mappingExo[15], // dialog: speed=
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

	baseExoJP = strings.ReplaceAll(baseExoJP, "\n", "\r\n")
	baseExoJP43 = strings.ReplaceAll(baseExoJP43, "\n", "\r\n")
	baseExoEN = strings.ReplaceAll(baseExoEN, "\n", "\r\n")
	baseExoEN43 = strings.ReplaceAll(baseExoEN43, "\n", "\r\n")
	baseExoJPv1 = strings.ReplaceAll(baseExoJPv1, "\n", "\r\n")
	baseExoJP43v1 = strings.ReplaceAll(baseExoJP43v1, "\n", "\r\n")
	baseExoENv1 = strings.ReplaceAll(baseExoENv1, "\n", "\r\n")
	baseExoEN43v1 = strings.ReplaceAll(baseExoEN43v1, "\n", "\r\n")

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
