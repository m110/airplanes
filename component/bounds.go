package component

import "github.com/yohamta/donburi"

type BoundsData struct {
	Disabled bool
}

var Bounds = donburi.NewComponentType[BoundsData]()

func GetBounds(entry *donburi.Entry) *BoundsData {
	return donburi.Get[BoundsData](entry, Bounds)
}
