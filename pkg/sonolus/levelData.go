package sonolus

type LevelData struct {
	BgmOffset float64           `json:"bgmOffset"`
	Entities  []LevelDataEntity `json:"entities"`
}

type LevelDataEntity struct {
	Archetype string                 `json:"archetype"`
	Data      []LevelDataEntityValue `json:"data"`
	Name      string                 `json:"name"`
}

type LevelDataEntityValue struct {
	Name  string
	Value float64
	Ref   string
}
