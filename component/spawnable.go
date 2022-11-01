package component

import "github.com/yohamta/donburi"

type SpawnFunc func(w donburi.World)

type SpawnableData struct {
	SpawnFunc SpawnFunc
}

var Spawnable = donburi.NewComponentType[SpawnableData]()

func GetSpawnable(entry *donburi.Entry) *SpawnableData {
	return donburi.Get[SpawnableData](entry, Spawnable)
}
