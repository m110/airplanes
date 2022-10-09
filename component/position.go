package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PositionData struct {
	Position engine.Vector
}

var Position = donburi.NewComponentType[PositionData]()

func GetPosition(entry *donburi.Entry) *PositionData {
	return donburi.Get[PositionData](entry, Position)
}
