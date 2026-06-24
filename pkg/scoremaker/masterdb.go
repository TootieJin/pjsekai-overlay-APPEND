package scoremaker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MusicMetadata struct {
	Title                string
	Lyricist             string
	Composer             string
	Arranger             string
	MusicCollaborationID int
	Label                string
}

type musicRow struct {
	ID                   int    `json:"id"`
	Title                string `json:"title"`
	Lyricist             string `json:"lyricist"`
	Composer             string `json:"composer"`
	Arranger             string `json:"arranger"`
	MusicCollaborationID int    `json:"musicCollaborationId"`
}

type musicCollaborationRow struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func dbBaseURL(region string) string {
	region = strings.ToLower(strings.TrimSpace(region))
	if region == "" || region == "jp" {
		return "https://raw.githubusercontent.com/Sekai-World/sekai-master-db-diff/refs/heads/main"
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/Sekai-World/sekai-master-db-%s-diff/refs/heads/main", region)
}

func FetchMusicMetadata(region string, musicID int, client *http.Client) (MusicMetadata, error) {
	if musicID <= 0 {
		return MusicMetadata{}, fmt.Errorf("music id must be positive")
	}
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}

	base := dbBaseURL(region)
	musicsURL := base + "/musics.json"
	collabsURL := base + "/musicCollaborations.json"

	req, err := http.NewRequest(http.MethodGet, musicsURL, nil)
	if err != nil {
		return MusicMetadata{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return MusicMetadata{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return MusicMetadata{}, fmt.Errorf("failed to fetch musics.json: status=%d", resp.StatusCode)
	}

	var musics []musicRow
	if err := json.NewDecoder(resp.Body).Decode(&musics); err != nil {
		return MusicMetadata{}, err
	}

	var target *musicRow
	for i := range musics {
		if musics[i].ID == musicID {
			target = &musics[i]
			break
		}
	}
	if target == nil {
		return MusicMetadata{}, fmt.Errorf("music id %d not found in musics.json", musicID)
	}

	meta := MusicMetadata{
		Title:                strings.TrimSpace(target.Title),
		Lyricist:             strings.TrimSpace(target.Lyricist),
		Composer:             strings.TrimSpace(target.Composer),
		Arranger:             strings.TrimSpace(target.Arranger),
		MusicCollaborationID: target.MusicCollaborationID,
	}

	if meta.MusicCollaborationID > 0 {
		req2, err := http.NewRequest(http.MethodGet, collabsURL, nil)
		if err != nil {
			return MusicMetadata{}, err
		}
		resp2, err := client.Do(req2)
		if err != nil {
			return MusicMetadata{}, err
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == http.StatusOK {
			var rows []musicCollaborationRow
			if err := json.NewDecoder(resp2.Body).Decode(&rows); err == nil {
				for _, row := range rows {
					if row.ID == meta.MusicCollaborationID {
						meta.Label = strings.TrimSpace(row.Label)
						break
					}
				}
			}
		}
	}

	return meta, nil
}
