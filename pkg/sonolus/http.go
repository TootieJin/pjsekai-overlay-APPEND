package sonolus

import (
	"net/url"
	"strings"
)

type LevelInfo struct {
	Name          string                  `json:"name"`
	Title         string                  `json:"title"`
	Artists       string                  `json:"artists"`
	Author        string                  `json:"author"`
	Version       int                     `json:"version"`
	Rating        int                     `json:"rating"`
	Cover         SRL                     `json:"cover"`
	Data          SRL                     `json:"data"`
	UseBackground UseItem[BackgroundInfo] `json:"useBackground"`
	Engine        EngineInfo              `json:"engine"`
	Tags          []TagInfo               `json:"tags"`
}

type TagInfo struct {
	Title string `json:"title"`
}

type BackgroundInfo struct {
	Image SRL `json:"image"`
}

type EngineInfo struct {
	Version    int            `json:"version"`
	Background BackgroundInfo `json:"background"`
}

type SRL struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

type InfoResponse[T any] struct {
	Item T `json:"item"`
}

type UseItem[T any] struct {
	UseDefault bool `json:"useDefault"`
	Item       T    `json:"item"`
}

func JoinUrl(base string, path string) (string, error) {
	if strings.HasPrefix(path, "http") {
		return path, nil
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path
	return u.String(), nil

}
