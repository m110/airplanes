package component

import "github.com/yohamta/donburi"

type BoundsData struct {
	Disabled bool
}

// Bounds indicates that the entity can't move of out the screen.
var Bounds = donburi.NewComponentType[BoundsData]()
