package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PositionData struct {
	Position engine.Vector
	Parent   *donburi.Entry
}

func (d *PositionData) SetParent(parent *donburi.Entry) {
	absPos := GetPosition(parent).AbsolutePosition()

	d.Parent = parent
	d.Position.X -= absPos.X
	d.Position.Y -= absPos.Y
}

func (d *PositionData) AbsolutePosition() engine.Vector {
	if d.Parent == nil {
		return d.Position
	}

	parent := GetPosition(d.Parent)
	return parent.AbsolutePosition().Add(d.Position)
}

var Position = donburi.NewComponentType[PositionData]()

func GetPosition(entry *donburi.Entry) *PositionData {
	return donburi.Get[PositionData](entry, Position)
}
