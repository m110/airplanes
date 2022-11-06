package component

import "github.com/yohamta/donburi"

type DespawnableData struct {
	// Set when the unit first appears on-screen
	Spawned bool
}

var Despawnable = donburi.NewComponentType[DespawnableData]()
