package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type TransformData struct {
	Position engine.Vector
	Rotation float64

	Parent   *donburi.Entry
	Children []*donburi.Entry
}

func (d *TransformData) AppendChild(parent, child *donburi.Entry) {
	d.Children = append(d.Children, child)
	GetTransform(child).SetParent(parent)
}

func (d *TransformData) SetParent(parent *donburi.Entry) {
	absPos := GetTransform(parent).WorldPosition()

	d.Parent = parent
	d.Position.X -= absPos.X
	d.Position.Y -= absPos.Y
}

func (d *TransformData) WorldPosition() engine.Vector {
	if d.Parent == nil {
		return d.Position
	}

	parent := GetTransform(d.Parent)
	return parent.WorldPosition().Add(d.Position)
}

var Transform = donburi.NewComponentType[TransformData]()

func GetTransform(entry *donburi.Entry) *TransformData {
	return donburi.Get[TransformData](entry, Transform)
}
