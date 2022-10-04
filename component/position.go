package component

import "github.com/yohamta/donburi"

type PositionData struct {
	X float64
	Y float64
}

var Position = donburi.NewComponentType[PositionData]()

func GetPosition(entry *donburi.Entry) *PositionData {
	return donburi.Get[PositionData](entry, Position)
}
