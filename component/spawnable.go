package component

import "github.com/yohamta/donburi"

type SpawnFunc func(w donburi.World)

type SpawnableData struct {
	SpawnFunc SpawnFunc
}

var Spawnable = donburi.NewComponentType[SpawnableData]()
