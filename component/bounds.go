package component

import "github.com/yohamta/donburi"

type BoundsData struct {
	Disabled bool
}

// Bounds indicates that the entity can't move of out the screen.
var Bounds = donburi.NewComponentType[BoundsData]()

func GetBounds(entry *donburi.Entry) *BoundsData {
	return donburi.Get[BoundsData](entry, Bounds)
}
