package component

import (
	"math"

	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type TransformData struct {
	LocalPosition engine.Vector
	LocalRotation float64
	LocalScale    engine.Vector

	Parent   *donburi.Entry
	Children []*donburi.Entry
}

func (d *TransformData) AppendChild(parent, child *donburi.Entry, keepWorldPosition bool) {
	d.Children = append(d.Children, child)
	GetTransform(child).SetParent(parent, keepWorldPosition)
}

func (d *TransformData) SetParent(parent *donburi.Entry, keepWorldPosition bool) {
	d.Parent = parent
	if keepWorldPosition {
		parentTransform := GetTransform(parent)
		parentPos := parentTransform.WorldPosition()

		d.LocalPosition = d.LocalPosition.Sub(parentPos)
		d.LocalRotation -= parentTransform.WorldRotation()
		d.LocalScale = d.LocalScale.Sub(parentTransform.WorldScale())
	}
}

func (d *TransformData) SetWorldPosition(pos engine.Vector) {
	if d.Parent == nil {
		d.LocalPosition = pos
		return
	}

	parentPos := GetTransform(d.Parent).WorldPosition()
	d.LocalPosition.X = pos.X - parentPos.X
	d.LocalPosition.Y = pos.Y - parentPos.Y
}

func (d *TransformData) WorldPosition() engine.Vector {
	if d.Parent == nil {
		return d.LocalPosition
	}

	parent := GetTransform(d.Parent)
	return parent.WorldPosition().Add(d.LocalPosition)
}

func (d *TransformData) SetWorldRotation(rotation float64) {
	if d.Parent == nil {
		d.LocalRotation = rotation
		return
	}

	parent := GetTransform(d.Parent)
	d.LocalRotation = rotation - parent.WorldRotation()
}

func (d *TransformData) WorldRotation() float64 {
	if d.Parent == nil {
		return d.LocalRotation
	}

	parent := GetTransform(d.Parent)
	return parent.WorldRotation() + d.LocalRotation
}

func (d *TransformData) WorldScale() engine.Vector {
	if d.Parent == nil {
		return d.LocalScale
	}

	parent := GetTransform(d.Parent)
	return parent.WorldScale().Add(d.LocalScale)
}

func (d *TransformData) Right() engine.Vector {
	radians := engine.ToRadians(d.WorldRotation())
	return engine.Vector{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

func (d *TransformData) Up() engine.Vector {
	radians := engine.ToRadians(d.WorldRotation() - 90.0)
	return engine.Vector{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

func (d *TransformData) LookAt(target engine.Vector) {
	x := target.X - d.WorldPosition().X
	y := target.Y - d.WorldPosition().Y
	radians := math.Atan2(y, x)
	d.SetWorldRotation(engine.ToDegrees(radians))
}

var Transform = donburi.NewComponentType[TransformData]()

func GetTransform(entry *donburi.Entry) *TransformData {
	return donburi.Get[TransformData](entry, Transform)
}
