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
	Time      float64
	SkillTime float64
	Score     float64
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

	"Skill":       0,
	"FeverChance": 0,
	"FeverStart":  0,
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

func CalculateScore(levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, power float64, scoreMode string) []PedFrame {
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
	frames = append(frames, PedFrame{
		Time:      0,
		SkillTime: 0,
		Score:     0,
	})
	bpmChanges := ([]BpmChange{})
	levelFax := float64(rating-5)*0.005 + 1
	comboFax := 1.0
	skillFax := 1.0
	comboCounterFax := 1.0

	score := 0.0
	entityCounter := 0
	comboCounter := 0
	noteEntities := ([]sonolus.LevelDataEntity{})

	var weightedComboFax float64
	if scoreMode == "sonolus" {
		for _, entity := range levelData.Entities {
			weight := WEIGHT_MAP[entity.Archetype]
			if weight == 0 {
				continue
			}
			weight = 1

			comboCounter += int(weight)
			if comboCounter%100 == 1 && comboCounter > 1 {
				comboCounterFax += 0.01
			}
			if comboCounterFax > 1.1 {
				comboCounterFax = 1.1
			}

			weightedComboFax += comboCounterFax
		}
	}

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
		} else if entity.Archetype == "Skill" {
			noteEntities = append(noteEntities, entity)
		}
	}
	slices.SortStableFunc(noteEntities, func(a, b sonolus.LevelDataEntity) int {
		aBeat, _ := getValueFromData(a.Data, "#BEAT")
		bBeat, _ := getValueFromData(b.Data, "#BEAT")
		if aBeat < bBeat {
			return -1
		}
		if aBeat > bBeat {
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

	skillActiveUntil := 0.0
	for _, entity := range noteEntities {
		weight := WEIGHT_MAP[entity.Archetype]

		// get beat/time early so skillFax can be applied to this event
		beat, err := getValueFromData(entity.Data, "#BEAT")
		if err != nil {
			continue
		}
		eventTime := getTimeFromBpmChanges(bpmChanges, beat)

		entityCounter += 1
		if entityCounter%100 == 1 && entityCounter > 1 {
			comboFax += 0.01
		}
		if comboFax > 1.1 {
			comboFax = 1.1
		}

		if eventTime <= skillActiveUntil {
			skillFax = 2.0
		} else {
			skillFax = 1.0
		}

		if entity.Archetype == "Skill" {
			// activate skill for 5 seconds starting at eventTime
			skillFax = 2.0
			skillActiveUntil = eventTime + 5.0

			frames = append(frames, PedFrame{
				Time:      eventTime,
				SkillTime: eventTime,
				Score:     score,
			})
			continue
		}

		switch scoreMode {
		default:
			score += ((float64(power) / weightedNotesCount) * // Team power / weighted notes count
				4 * // Constant
				weight * // Note weight
				1 * // Judge weight (Always 1)
				levelFax * // Level fax
				comboFax * // Combo fax
				skillFax) // Skill fax
		case "tournament":
			score += 3
		case "sonolus":
			score += 1000000 * (comboFax / weightedComboFax)
		}

		frames = append(frames, PedFrame{
			Time:      eventTime,
			SkillTime: 0,
			Score:     score,
		})
	}

	return frames
}

func WritePedFile(frames []PedFrame, assets string, path string, levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, scoreMode string, enUI bool) error {
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
		if scoreMode == "tournament" {
			weight = 3
		}
		weightedNotesCount += weight
	}

	lastScore := 0.0
	rating := math.Max(5, math.Min(float64(levelInfo.Rating), 40))
	procCount := 0
	for i, frame := range frames {
		score := frame.Score
		frameScore := math.Trunc(score - lastScore)
		lastScore = frame.Score

		rank := "n"
		scoreX := 0.0
		scoreXv1 := 0.0

		rankBorder := float64(1200000 + (rating-5)*4100)
		rankS := float64(1040000 + (rating-5)*5200)
		rankA := float64(840000 + (rating-5)*4200)
		rankB := float64(400000 + (rating-5)*2000)
		rankC := float64(20000 + (rating-5)*100)

		switch scoreMode {
		case "tournament":
			rankBorder = math.Floor(weightedNotesCount)
			rankS = math.Floor(weightedNotesCount * 0.800)
			rankA = math.Floor(weightedNotesCount * 0.666)
			rankB = math.Floor(weightedNotesCount * 0.533)
			rankC = math.Floor(weightedNotesCount * 0.400)
		case "sonolus":
			rankBorder = 1000000
			rankS = 900000
			rankA = 750000
			rankB = 600000
			rankC = 450000
		}

		// bar

		rankBorderPos := 1650.0 / 1650
		rankSPos := 1478.0 / 1650
		rankAPos := 1234.0 / 1650
		rankBPos := 990.0 / 1650
		rankCPos := 746.0 / 1650

		rankBorderPosv1 := 610.0 / 610
		rankSPosv1 := 541.5 / 610
		rankAPosv1 := 452.0 / 610
		rankBPosv1 := 361.0 / 610
		rankCPosv1 := 273.0 / 610

		if score < 0 {
			rank = "d"
			scoreX = 0
			scoreXv1 = 0
		} else if score >= rankBorder {
			rank = "s"
			scoreX = rankBorderPos
			scoreXv1 = rankBorderPosv1
		} else if score >= rankS {
			rank = "s"
			scoreX = (float64((score-rankS))/float64((rankBorder-rankS)))*(rankBorderPos-rankSPos) + rankSPos
			scoreXv1 = (float64((score-rankS))/float64((rankBorder-rankS)))*(rankBorderPosv1-rankSPosv1) + rankSPosv1
		} else if score >= rankA {
			rank = "a"
			scoreX = (float64((score-rankA))/float64((rankS-rankA)))*(rankSPos-rankAPos) + rankAPos
			scoreXv1 = (float64((score-rankA))/float64((rankS-rankA)))*(rankSPosv1-rankAPosv1) + rankAPosv1
		} else if score >= rankB {
			rank = "b"
			scoreX = (float64((score-rankB))/float64((rankA-rankB)))*(rankAPos-rankBPos) + rankBPos
			scoreXv1 = (float64((score-rankB))/float64((rankA-rankB)))*(rankAPosv1-rankBPosv1) + rankBPosv1
		} else if score >= rankC {
			rank = "c"
			scoreX = (float64((score-rankC))/float64((rankB-rankC)))*(rankBPos-rankCPos) + rankCPos
			scoreXv1 = (float64((score-rankC))/float64((rankB-rankC)))*(rankBPosv1-rankCPosv1) + rankCPosv1
		} else {
			rank = "d"
			scoreX = (float64(score) / float64(rankC)) * rankCPos
			scoreXv1 = (float64(score) / float64(rankC)) * rankCPosv1
		}

		time := frame.Time
		skillTime := frame.SkillTime

		if skillTime > 0 && skillTime == time {
			// only emit when this frame is the actual activation frame and avoid duplicates
			if i == 0 || frames[i-1].SkillTime != skillTime {
				writer.Write(fmt.Appendf(nil, "s|%f\n", skillTime))
				procCount++
			}
		}

		combo := i - procCount

		if combo%100 == 0 && combo > 0 {
			writer.Write(fmt.Appendf(nil, "c|%f:%d\n", time, combo))
		}
		if i < len(frames)-1 && time == frames[i+1].Time {
			continue
		}
		if time == 0 && i > 0 {
			time = frames[i-1].Time + 0.000001
		}

		writer.Write(fmt.Appendf(nil, "d|%f:%.0f:%.0f:%f:%f:%s:%d\n", time, score, frameScore, scoreX, scoreXv1, rank, combo))
	}

	return nil
}
