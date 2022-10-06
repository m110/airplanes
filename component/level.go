package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/engine"
)

type LevelData struct {
	ProgressionTimer *engine.Timer
	ReachedEnd       bool
	Progressed       bool
}

var Level = donburi.NewComponentType[LevelData]()

func GetLevel(entry *donburi.Entry) *LevelData {
	return donburi.Get[LevelData](entry, Level)
}

func MustFindLevel(w donburi.World) *donburi.Entry {
	level, ok := query.NewQuery(filter.Contains(Level)).FirstEntity(w)
	if !ok {
		panic("no level found")
	}

	return level
}
