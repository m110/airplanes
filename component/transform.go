package component

import (
	"math"

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

func (d TransformData) WorldRotation() float64 {
	if d.Parent == nil {
		return d.Rotation
	}

	parent := GetTransform(d.Parent)
	return parent.WorldRotation() + d.Rotation
}

func (d *TransformData) Right() engine.Vector {
	radians := engine.ToRadians(d.Rotation)
	return engine.Vector{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

func (d *TransformData) Up() engine.Vector {
	radians := engine.ToRadians(d.Rotation - 90.0)
	return engine.Vector{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

var Transform = donburi.NewComponentType[TransformData]()

func GetTransform(entry *donburi.Entry) *TransformData {
	return donburi.Get[TransformData](entry, Transform)
}
