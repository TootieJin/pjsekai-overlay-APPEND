package scoremaker

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
)

//go:embed dependencies/sonolus_converters.zip
var embeddedConvertersZip []byte

var (
	embeddedConverterRootOnce sync.Once
	embeddedConverterRoot     string
	embeddedConverterRootErr  error
)

func DecodePublishedScoreInfo(raw []byte) (PublishedScoreInfo, error) {
	return extractPublishedScoreInfo(raw)
}

func LoadPublishedScoreInfoFile(path string) (PublishedScoreInfo, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return PublishedScoreInfo{}, err
	}
	return DecodePublishedScoreInfo(raw)
}

func isValidConverterRoot(root string) bool {
	if strings.TrimSpace(root) == "" {
		return false
	}
	stat, err := os.Stat(filepath.Join(root, "sonolus_converters"))
	return err == nil && stat.IsDir()
}

func ensureEmbeddedConverterRoot() (string, error) {
	embeddedConverterRootOnce.Do(func() {
		if len(embeddedConvertersZip) == 0 {
			embeddedConverterRootErr = fmt.Errorf("embedded sonolus_converters archive is empty")
			return
		}

		cacheBase, err := os.UserCacheDir()
		if err != nil || strings.TrimSpace(cacheBase) == "" {
			cacheBase = os.TempDir()
		}
		root := filepath.Join(cacheBase, "pjsekai-overlay-APPEND", "sonolus-level-converters")
		if isValidConverterRoot(root) {
			embeddedConverterRoot = root
			return
		}

		if err := os.MkdirAll(root, 0755); err != nil {
			embeddedConverterRootErr = err
			return
		}

		reader, err := zip.NewReader(bytes.NewReader(embeddedConvertersZip), int64(len(embeddedConvertersZip)))
		if err != nil {
			embeddedConverterRootErr = err
			return
		}

		for _, file := range reader.File {
			targetPath := filepath.Join(root, file.Name)
			cleanRoot := filepath.Clean(root) + string(os.PathSeparator)
			cleanTarget := filepath.Clean(targetPath)
			if !strings.HasPrefix(cleanTarget, cleanRoot) && cleanTarget != filepath.Clean(root) {
				embeddedConverterRootErr = fmt.Errorf("embedded archive contains invalid path: %s", file.Name)
				return
			}

			if file.FileInfo().IsDir() {
				if err := os.MkdirAll(cleanTarget, 0755); err != nil {
					embeddedConverterRootErr = err
					return
				}
				continue
			}

			if err := os.MkdirAll(filepath.Dir(cleanTarget), 0755); err != nil {
				embeddedConverterRootErr = err
				return
			}

			src, err := file.Open()
			if err != nil {
				embeddedConverterRootErr = err
				return
			}

			dst, err := os.Create(cleanTarget)
			if err != nil {
				src.Close()
				embeddedConverterRootErr = err
				return
			}

			_, copyErr := io.Copy(dst, src)
			closeErr := dst.Close()
			src.Close()
			if copyErr != nil {
				embeddedConverterRootErr = copyErr
				return
			}
			if closeErr != nil {
				embeddedConverterRootErr = closeErr
				return
			}
		}

		if !isValidConverterRoot(root) {
			embeddedConverterRootErr = fmt.Errorf("embedded sonolus_converters extraction did not produce a valid root")
			return
		}
		embeddedConverterRoot = root
	})

	if embeddedConverterRootErr != nil {
		return "", embeddedConverterRootErr
	}
	return embeddedConverterRoot, nil
}

func findConverterRoot() (string, error) {
	if env := strings.TrimSpace(os.Getenv("SONOLUS_CONVERTERS_DIR")); env != "" {
		if isValidConverterRoot(env) {
			return env, nil
		}
		return "", fmt.Errorf("SONOLUS_CONVERTERS_DIR does not contain a valid sonolus_converters package: %s", env)
	}

	embeddedRoot, err := ensureEmbeddedConverterRoot()
	if err == nil {
		return embeddedRoot, nil
	}

	return "", fmt.Errorf("embedded sonolus_converters setup failed: %w", err)
}

func maybeUnwrapOuterGzipBase64(raw []byte) []byte {
	tryBase64 := func(data []byte) ([]byte, bool) {
		trimmed := strings.TrimSpace(string(data))
		if trimmed == "" {
			return nil, false
		}
		decoded, err := base64.StdEncoding.DecodeString(trimmed)
		if err != nil {
			return nil, false
		}
		if len(decoded) >= 2 && decoded[0] == 0x1f && decoded[1] == 0x8b {
			return []byte(trimmed), true
		}
		return nil, false
	}

	if b64, ok := tryBase64(raw); ok {
		return b64
	}

	// Common downloaded format: gzip(base64(gzip(json)))
	if len(raw) >= 2 && raw[0] == 0x1f && raw[1] == 0x8b {
		gz, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return raw
		}
		defer gz.Close()

		inner, err := io.ReadAll(gz)
		if err != nil {
			return raw
		}
		if b64, ok := tryBase64(inner); ok {
			return b64
		}
	}

	return raw
}

func runPJSKConverter(raw []byte, exportCode string) ([]byte, error) {
	converterRoot, err := findConverterRoot()
	if err != nil {
		return nil, err
	}

	raw = maybeUnwrapOuterGzipBase64(raw)

	pythonBin := "python"
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("python"); err != nil {
			if _, fallbackErr := exec.LookPath("py"); fallbackErr == nil {
				pythonBin = "py"
			} else {
				return nil, fmt.Errorf("python interpreter not found: %w", err)
			}
		}
	} else if _, err := exec.LookPath(pythonBin); err != nil {
		return nil, err
	}

	converterCode := `import base64
import gzip
import io
import json
import os
import sys

root = os.environ.get("SONOLUS_CONVERTERS_DIR")
if not root:
    raise SystemExit("SONOLUS_CONVERTERS_DIR is required")
sys.path.insert(0, root)

from sonolus_converters import LevelData, pjsk, sus, usc

raw = base64.b64decode(sys.stdin.buffer.read())

def maybe_b64(data: bytes) -> bytes | None:
	try:
		text = data.decode("ascii").strip()
	except Exception:
		return None
	if not text:
		return None
	try:
		return base64.b64decode(text, validate=False)
	except Exception:
		return None

def maybe_gzip(data: bytes) -> bytes | None:
	if len(data) < 2 or data[0] != 0x1F or data[1] != 0x8B:
		return None
	try:
		return gzip.decompress(data)
	except Exception:
		return None

def build_candidates(seed: bytes, max_depth: int = 4) -> list[bytes]:
	out = [seed]
	seen = {seed}
	idx = 0
	depth = 0
	while idx < len(out) and depth < max_depth:
		cur = out[idx]
		idx += 1
		for nxt in (maybe_gzip(cur), maybe_b64(cur)):
			if nxt is not None and nxt not in seen:
				seen.add(nxt)
				out.append(nxt)
		depth += 1
	return out

score = None
selected_raw = None
last_err = None
for candidate in build_candidates(raw):
	try:
		score = pjsk.load(io.BytesIO(candidate))
		selected_raw = candidate
		break
	except Exception as exc:
		last_err = exc

if score is None:
	raise last_err if last_err is not None else RuntimeError("unable to decode chart bytes")

` + exportCode

	cmd := exec.Command(pythonBin, "-c", converterCode)
	cmd.Env = append(os.Environ(), "SONOLUS_CONVERTERS_DIR="+converterRoot)
	cmd.Stdin = bytes.NewReader([]byte(base64.StdEncoding.EncodeToString(raw)))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("converter failed: %w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return nil, err
	}

	return stdout.Bytes(), nil
}

type susEventExtract struct {
	Archetype string  `json:"archetype"`
	Beat      float64 `json:"beat"`
}

type susLevelDataExtract struct {
	LevelData sonolus.LevelData `json:"levelData"`
	Events    []susEventExtract `json:"events"`
}

func AppendEventArchetypes(levelData *sonolus.LevelData, events []susEventExtract) {
	if levelData == nil || len(events) == 0 {
		return
	}
	exists := make(map[string]struct{})
	for _, entity := range levelData.Entities {
		if entity.Archetype != "Skill" && entity.Archetype != "FeverChance" && entity.Archetype != "FeverStart" {
			continue
		}
		beat := 0.0
		found := false
		for _, value := range entity.Data {
			if value.Name == "#BEAT" {
				beat = value.Value
				found = true
				break
			}
		}
		if !found {
			continue
		}
		key := fmt.Sprintf("%s:%.6f", entity.Archetype, beat)
		exists[key] = struct{}{}
	}

	for _, event := range events {
		if event.Archetype != "Skill" && event.Archetype != "FeverChance" && event.Archetype != "FeverStart" {
			continue
		}
		beat := math.Round(event.Beat*1e6) / 1e6
		key := fmt.Sprintf("%s:%.6f", event.Archetype, beat)
		if _, ok := exists[key]; ok {
			continue
		}
		levelData.Entities = append(levelData.Entities, sonolus.LevelDataEntity{
			Archetype: event.Archetype,
			Data: []sonolus.LevelDataEntityValue{
				{Name: "#BEAT", Value: beat},
			},
			Name: "",
		})
		exists[key] = struct{}{}
	}
}

func AppendEventArchetypesFromLevelData(dst *sonolus.LevelData, src sonolus.LevelData) {
	if dst == nil || len(src.Entities) == 0 {
		return
	}
	events := make([]susEventExtract, 0)
	for _, entity := range src.Entities {
		if entity.Archetype != "Skill" && entity.Archetype != "FeverChance" && entity.Archetype != "FeverStart" {
			continue
		}
		beat := 0.0
		for _, value := range entity.Data {
			if value.Name == "#BEAT" {
				beat = value.Value
				break
			}
		}
		events = append(events, susEventExtract{Archetype: entity.Archetype, Beat: beat})
	}
	AppendEventArchetypes(dst, events)
}

func convertSUSTextToNextSekaiLevelData(raw []byte) (sonolus.LevelData, error) {
	converterRoot, err := findConverterRoot()
	if err != nil {
		return sonolus.LevelData{}, err
	}

	pythonBin := "python"
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("python"); err != nil {
			if _, fallbackErr := exec.LookPath("py"); fallbackErr == nil {
				pythonBin = "py"
			} else {
				return sonolus.LevelData{}, fmt.Errorf("python interpreter not found: %w", err)
			}
		}
	} else if _, err := exec.LookPath(pythonBin); err != nil {
		return sonolus.LevelData{}, err
	}

	converterCode := `import io
import json
import os
import sys

root = os.environ.get("SONOLUS_CONVERTERS_DIR")
if not root:
    raise SystemExit("SONOLUS_CONVERTERS_DIR is required")
sys.path.insert(0, root)

from sonolus_converters import LevelData, sus
from sonolus_converters.notes.single import Skill, FeverChance, FeverStart

sus_text = sys.stdin.buffer.read().decode("utf-8", errors="ignore")
score = sus.load(io.StringIO(sus_text))

buf = io.BytesIO()
LevelData.next_sekai.export(buf, score, as_compressed=False)
level_data = json.loads(buf.getvalue().decode("utf-8"))

events = []
for note in score.notes:
    if isinstance(note, Skill):
        events.append({"archetype": "Skill", "beat": float(note.beat)})
    elif isinstance(note, FeverChance):
        events.append({"archetype": "FeverChance", "beat": float(note.beat)})
    elif isinstance(note, FeverStart):
        events.append({"archetype": "FeverStart", "beat": float(note.beat)})

sys.stdout.buffer.write(json.dumps({"levelData": level_data, "events": events}).encode("utf-8"))
`

	cmd := exec.Command(pythonBin, "-c", converterCode)
	cmd.Env = append(os.Environ(), "SONOLUS_CONVERTERS_DIR="+converterRoot)
	cmd.Stdin = bytes.NewReader(raw)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return sonolus.LevelData{}, fmt.Errorf("sus converter failed: %w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return sonolus.LevelData{}, err
	}

	var extracted susLevelDataExtract
	if err := json.Unmarshal(stdout.Bytes(), &extracted); err != nil {
		return sonolus.LevelData{}, err
	}
	AppendEventArchetypes(&extracted.LevelData, extracted.Events)
	return extracted.LevelData, nil
}

func ConvertSUSURLToNextSekaiLevelData(susURL string) (sonolus.LevelData, error) {
	if strings.TrimSpace(susURL) == "" {
		return sonolus.LevelData{}, fmt.Errorf("sus url is required")
	}

	resp, err := http.Get(susURL)
	if err != nil {
		return sonolus.LevelData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return sonolus.LevelData{}, fmt.Errorf("failed to fetch sus file: status=%d body=%q", resp.StatusCode, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return sonolus.LevelData{}, err
	}
	return convertSUSTextToNextSekaiLevelData(raw)
}

func ConvertSUSMusicIDToNextSekaiLevelData(region string, musicID int) (sonolus.LevelData, error) {
	if musicID <= 0 {
		return sonolus.LevelData{}, fmt.Errorf("music id must be > 0")
	}
	susURL := fmt.Sprintf("https://storage.sekai.best/sekai-%s-assets/music/music_score/%04d_01/master.txt", region, musicID)
	levelData, err := ConvertSUSURLToNextSekaiLevelData(susURL)
	if err == nil || strings.EqualFold(strings.TrimSpace(region), "jp") {
		return levelData, err
	}
	fallbackURL := fmt.Sprintf("https://storage.sekai.best/sekai-jp-assets/music/music_score/%04d_01/master.txt", musicID)
	fallbackLevelData, fallbackErr := ConvertSUSURLToNextSekaiLevelData(fallbackURL)
	if fallbackErr == nil {
		return fallbackLevelData, nil
	}
	return sonolus.LevelData{}, fmt.Errorf("failed to fetch sus file for region %s and jp fallback: %w; fallback: %v", region, err, fallbackErr)
}

func ConvertPJSKToNextSekaiLevelData(raw []byte) (sonolus.LevelData, error) {
	stdout, err := runPJSKConverter(raw, `buf = io.BytesIO()
LevelData.next_sekai.export(buf, score, as_compressed=False)
sys.stdout.buffer.write(buf.getvalue())
`)
	if err != nil {
		return sonolus.LevelData{}, fmt.Errorf("next sekai conversion failed: %w", err)
	}

	var levelData sonolus.LevelData
	if err := json.Unmarshal(stdout, &levelData); err != nil {
		return sonolus.LevelData{}, err
	}
	return levelData, nil
}

func ConvertPJSKToUSC(raw []byte) ([]byte, error) {
	stdout, err := runPJSKConverter(raw, `buf = io.BytesIO()
text = io.StringIO()
usc.export(text, score, minified=False)
sys.stdout.write(text.getvalue())
`)
	if err != nil {
		return nil, fmt.Errorf("usc conversion failed: %w", err)
	}
	return stdout, nil
}

func ConvertRawPJSKToReadablePJSK(raw []byte) ([]byte, error) {
	stdout, err := runPJSKConverter(raw, `pjsk_data = pjsk.load_raw(io.BytesIO(selected_raw))
sys.stdout.buffer.write(json.dumps(pjsk_data, indent=2, ensure_ascii=False).encode("utf-8"))
`)
	if err != nil {
		return nil, fmt.Errorf("readable .pjsk conversion failed: %w", err)
	}
	return stdout, nil
}

func ConvertPJSKToMMWS(raw []byte) ([]byte, error) {
	return ConvertPJSKToUSC(raw)
}
