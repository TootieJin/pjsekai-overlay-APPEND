package pjsekaioverlay

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/sonolus"
)

type PedFrame struct {
	Time  float64
	Score float64
}

type BpmChange struct {
	Beat float64
	Bpm  float64
}

var WEIGHT_MAP = map[string]float64{
	// Nanashi's Archetype
	"#BPM_CHANGE":    0,
	"Initialization": 0,
	"InputManager":   0,
	"Stage":          0,

	"NormalTapNote":   1,
	"CriticalTapNote": 2,

	"NormalFlickNote":   1,
	"CriticalFlickNote": 3,

	"NormalSlideStartNote":   1,
	"CriticalSlideStartNote": 2,

	"NormalSlideEndNote":   1,
	"CriticalSlideEndNote": 2,

	"NormalSlideEndFlickNote":   1,
	"CriticalSlideEndFlickNote": 3,

	"HiddenSlideTickNote":   0,
	"NormalSlideTickNote":   0.1,
	"CriticalSlideTickNote": 0.2,

	"IgnoredSlideTickNote":          0.1,
	"NormalAttachedSlideTickNote":   0.1,
	"CriticalAttachedSlideTickNote": 0.2,

	"NormalSlideConnector":   0,
	"CriticalSlideConnector": 0,

	"SimLine": 0,

	"NormalSlotEffect":       0,
	"SlideSlotEffect":        0,
	"FlickSlotEffect":        0,
	"CriticalSlotEffect":     0,
	"NormalSlotGlowEffect":   0,
	"SlideSlotGlowEffect":    0,
	"FlickSlotGlowEffect":    0,
	"CriticalSlotGlowEffect": 0,

	"NormalTraceNote":   0.1,
	"CriticalTraceNote": 0.2,

	"NormalTraceSlotEffect":     0,
	"NormalTraceSlotGlowEffect": 0,

	"DamageNote":           0.1,
	"DamageSlotEffect":     0,
	"DamageSlotGlowEffect": 0,

	"NormalTraceFlickNote":         1,
	"CriticalTraceFlickNote":       3,
	"NonDirectionalTraceFlickNote": 1,

	"NormalTraceSlideStartNote":   0.1,
	"NormalTraceSlideEndNote":     0.1,
	"CriticalTraceSlideStartNote": 0.2,
	"CriticalTraceSlideEndNote":   0.2,

	"TimeScaleGroup":  0,
	"TimeScaleChange": 0,

	// Next SEKAI's Archetype
	"#TIMESCALE_CHANGE": 0,
	"#TIMESCALE_GROUP":  0,
	"_InputManager":     0,
	"SlideManager":      0,
	"Connector":         0,
	"SlotGlowEffect":    0,
	"SlotEffect":        0,

	"NormalHeadTapNote":          1,
	"CriticalHeadTapNote":        2,
	"NormalHeadFlickNote":        1,
	"CriticalHeadFlickNote":      3,
	"NormalHeadTraceNote":        0.1,
	"CriticalHeadTraceNote":      0.2,
	"NormalHeadTraceFlickNote":   1,
	"CriticalHeadTraceFlickNote": 3,
	"NormalHeadReleaseNote":      1,
	"CriticalHeadReleaseNote":    2,

	"NormalTailTapNote":          1,
	"CriticalTailTapNote":        2,
	"NormalTailFlickNote":        1,
	"CriticalTailFlickNote":      3,
	"NormalTailTraceNote":        0.1,
	"CriticalTailTraceNote":      0.2,
	"NormalTailTraceFlickNote":   1,
	"CriticalTailTraceFlickNote": 3,
	"NormalTailReleaseNote":      1,
	"CriticalTailReleaseNote":    2,

	"TransientHiddenTickNote": 0.1,
	"NormalTickNote":          0.1,
	"CriticalTickNote":        0.2,
	"AnchorNote":              0,

	"FakeNormalTapNote":          0,
	"FakeCriticalTapNote":        0,
	"FakeNormalFlickNote":        0,
	"FakeCriticalFlickNote":      0,
	"FakeNormalTraceNote":        0,
	"FakeCriticalTraceNote":      0,
	"FakeNormalTraceFlickNote":   0,
	"FakeCriticalTraceFlickNote": 0,
	"FakeNormalReleaseNote":      0,
	"FakeCriticalReleaseNote":    0,

	"FakeNormalHeadTapNote":          0,
	"FakeCriticalHeadTapNote":        0,
	"FakeNormalHeadFlickNote":        0,
	"FakeCriticalHeadFlickNote":      0,
	"FakeNormalHeadTraceNote":        0,
	"FakeCriticalHeadTraceNote":      0,
	"FakeNormalHeadTraceFlickNote":   0,
	"FakeCriticalHeadTraceFlickNote": 0,
	"FakeNormalHeadReleaseNote":      0,
	"FakeCriticalHeadReleaseNote":    0,

	"FakeNormalTailTapNote":          0,
	"FakeCriticalTailTapNote":        0,
	"FakeNormalTailFlickNote":        0,
	"FakeCriticalTailFlickNote":      0,
	"FakeNormalTailTraceNote":        0,
	"FakeCriticalTailTraceNote":      0,
	"FakeNormalTailTraceFlickNote":   0,
	"FakeCriticalTailTraceFlickNote": 0,
	"FakeNormalTailReleaseNote":      0,
	"FakeCriticalTailReleaseNote":    0,

	"FakeTransientHiddenTickNote": 0,
	"FakeNormalTickNote":          0,
	"FakeCriticalTickNote":        0,
	"FakeAnchorNote":              0,

	"FakeDamageNote": 0,
}

func getValueFromData(data []sonolus.LevelDataEntityValue, name string) (float64, error) {
	for _, value := range data {
		if value.Name == name {
			return value.Value, nil
		}
	}
	return 0, fmt.Errorf("value not found: %s", name)
}

func getTimeFromBpmChanges(bpmChanges []BpmChange, beat float64) float64 {
	ret := 0.0
	for i, bpmChange := range bpmChanges {
		if i == len(bpmChanges)-1 {
			ret += (beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
			break
		}
		nextBpmChange := bpmChanges[i+1]
		if beat >= bpmChange.Beat && beat < nextBpmChange.Beat {
			ret += (beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
			break
		} else if beat >= nextBpmChange.Beat {
			ret += (nextBpmChange.Beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
		} else {
			break
		}
	}
	return ret
}

func CalculateScore(levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, power float64, exScore bool) []PedFrame {
	rating := levelInfo.Rating
	var weightedNotesCount float64 = 0
	for _, entity := range levelData.Entities {
		weight := WEIGHT_MAP[entity.Archetype]
		if weight == 0 {
			continue
		}
		weightedNotesCount += weight
	}

	frames := make([]PedFrame, 0, int(weightedNotesCount)+1)
	frames = append(frames, PedFrame{Time: 0, Score: 0})
	bpmChanges := ([]BpmChange{})
	levelFax := float64(rating-5)*0.005 + 1
	comboFax := 1.0

	score := 0.0
	entityCounter := 0
	noteEntities := ([]sonolus.LevelDataEntity{})

	for _, entity := range levelData.Entities {
		weight := WEIGHT_MAP[entity.Archetype]
		if weight > 0.0 && len(entity.Data) > 0 {
			noteEntities = append(noteEntities, entity)
		} else if entity.Archetype == "#BPM_CHANGE" {
			beat, err := getValueFromData(entity.Data, "#BEAT")
			if err != nil {
				continue
			}
			bpm, err := getValueFromData(entity.Data, "#BPM")
			if err != nil {
				continue
			}
			bpmChanges = append(bpmChanges, BpmChange{
				Beat: beat,
				Bpm:  bpm,
			})
		}
	}
	slices.SortStableFunc(noteEntities, func(a, b sonolus.LevelDataEntity) int {
		if a.Data[0].Value < b.Data[0].Value {
			return -1
		}
		if a.Data[0].Value > b.Data[0].Value {
			return 1
		}
		return 0
	})
	slices.SortStableFunc(bpmChanges, func(a, b BpmChange) int {
		if a.Beat < b.Beat {
			return -1
		}
		if a.Beat > b.Beat {
			return 1
		}
		return 0
	})
	for _, entity := range noteEntities {
		weight := WEIGHT_MAP[entity.Archetype]
		entityCounter += 1
		if entityCounter%100 == 1 && entityCounter > 1 {
			comboFax += 0.01
		}
		if comboFax > 1.1 {
			comboFax = 1.1
		}

		if exScore {
			score += 3
		} else {
			score += ((float64(power) / weightedNotesCount) * // Team power / weighted notes count
				4 * // Constant
				weight * // Note weight
				1 * // Judge weight (Always 1)
				levelFax * // Level fax
				comboFax * // Combo fax
				1) // Skill fax (Always 1)
		}

		beat, err := getValueFromData(entity.Data, "#BEAT")
		if err != nil {
			continue
		}
		frames = append(frames, PedFrame{
			Time:  getTimeFromBpmChanges(bpmChanges, beat) + levelData.BgmOffset,
			Score: score,
		})
	}

	return frames
}

func WritePedFile(frames []PedFrame, assets string, path string, levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, exScore bool, enUI bool) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました (Failed to create file.) [%s]", err)
	}
	defer file.Close()

	writer := io.Writer(file)

	fmt.Fprintf(writer, "p|%s\n", assets)
	fmt.Fprintf(writer, "e|%s\n", strconv.FormatBool(enUI))
	fmt.Fprintf(writer, "v|%s\n", Version)
	fmt.Fprintf(writer, "u|%d\n", time.Now().Unix())

	var weightedNotesCount float64 = 0
	for _, entity := range levelData.Entities {
		weight := WEIGHT_MAP[entity.Archetype]
		if weight == 0 {
			continue
		}
		if exScore {
			weight = 3
		}
		weightedNotesCount += weight
	}

	lastScore := 0.0
	lastScore2 := 0.0
	rating := levelInfo.Rating
	for i, frame := range frames {
		// 2-variable scoring (supports accurate digits up to 1e+34)
		score := math.Mod(frame.Score, 1e+17)
		score2 := math.Floor(frame.Score / 1e+17)
		if score2 < 0 {
			score2 = math.Ceil(frame.Score / 1e+17)
		}

		if math.Ceil(score) < 0 && math.Ceil(score2) < 0 {
			score = -score
		}

		frameScore := math.Mod(frame.Score-(lastScore+(lastScore2*1e+17)), 1e+17)
		frameScore2 := score2 - lastScore2

		lastScore = math.Mod(frame.Score, 1e+17)
		lastScore2 = math.Floor(frame.Score / 1e+17)
		if lastScore2 < 0 {
			lastScore2 = math.Ceil(frame.Score / 1e+17)
		}

		if math.Ceil(frameScore) < 0 && math.Floor(frameScore2) > 0 {
			frameScore = 1e+16 - frameScore
		}

		if math.Ceil(frameScore) < 0 && math.Ceil(frameScore2) < 0 {
			frameScore = -frameScore
		}

		rank := "n"
		scoreX := 0.0
		scoreXv1 := 0.0

		if rating < 5 {
			rating = 5
		} else if rating > 40 {
			rating = 40
		}

		rankBorder := float64(1200000 + (rating-5)*4100)
		rankS := float64(1040000 + (rating-5)*5200)
		rankA := float64(840000 + (rating-5)*4200)
		rankB := float64(400000 + (rating-5)*2000)
		rankC := float64(20000 + (rating-5)*100)

		if exScore {
			rankBorder = math.Floor(weightedNotesCount)
			rankS = math.Floor(weightedNotesCount * 0.800)
			rankA = math.Floor(weightedNotesCount * 0.666)
			rankB = math.Floor(weightedNotesCount * 0.533)
			rankC = math.Floor(weightedNotesCount * 0.400)
		}

		// bar
		if math.Ceil(score2) < 0 || math.Ceil(score) < 0 {
			rank = "d"
			scoreX = 0
			scoreXv1 = 0
		} else if score >= rankBorder || math.Floor(score2) > 0 {
			rank = "s"
			scoreX = 372
			scoreXv1 = 1
		} else if score >= rankS {
			rank = "s"
			scoreX = (float64((score-rankS))/float64((rankBorder-rankS)))*36 + 335
			scoreXv1 = (float64((score-rankS))/float64((rankBorder-rankS)))*0.110 + 0.890
		} else if score >= rankA {
			rank = "a"
			scoreX = (float64((score-rankA))/float64((rankS-rankA)))*55 + 280
			scoreXv1 = (float64((score-rankA))/float64((rankS-rankA)))*0.148 + 0.742
		} else if score >= rankB {
			rank = "b"
			scoreX = (float64((score-rankB))/float64((rankA-rankB)))*55 + 225
			scoreXv1 = (float64((score-rankB))/float64((rankA-rankB)))*0.151 + 0.591
		} else if score >= rankC {
			rank = "c"
			scoreX = (float64((score-rankC))/float64((rankB-rankC)))*55 + 170
			scoreXv1 = (float64((score-rankC))/float64((rankB-rankC)))*0.144 + 0.447
		} else {
			rank = "d"
			scoreX = (float64(score) / float64(rankC)) * 168
			scoreXv1 = (float64(score) / float64(rankC)) * 0.447
		}

		time := frame.Time
		if time == 0 && i > 0 {
			time = frames[i-1].Time + 0.000001
		}

		writer.Write(fmt.Appendf(nil, "s|%f:%.0f:%.0f:%.0f:%.0f:%f:%f:%s:%d\n", time, score2, score, frameScore2, frameScore, scoreX/372, scoreXv1, rank, i))
	}

	return nil
}
