package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PositionData struct {
	Position engine.Vector

	Parent   *donburi.Entry
	Children []*donburi.Entry
}

func (d *PositionData) AppendChild(parent, child *donburi.Entry) {
	d.Children = append(d.Children, child)
	GetPosition(child).SetParent(parent)
}

func (d *PositionData) SetParent(parent *donburi.Entry) {
	absPos := GetPosition(parent).WorldPosition()

	d.Parent = parent
	d.Position.X -= absPos.X
	d.Position.Y -= absPos.Y
}

func (d *PositionData) WorldPosition() engine.Vector {
	if d.Parent == nil {
		return d.Position
	}

	parent := GetPosition(d.Parent)
	return parent.WorldPosition().Add(d.Position)
}

var Position = donburi.NewComponentType[PositionData]()

func GetPosition(entry *donburi.Entry) *PositionData {
	return donburi.Get[PositionData](entry, Position)
}
