package archetype

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewAirbaseAirplane(w donburi.World, position math.Vec2, faction component.PlayerFaction, index int) {
	airplane := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
			component.Velocity,
			component.PlayerSelect,
			component.Altitude,
		),
	)

	originalRotation := -90.0

	donburi.SetValue(airplane, transform.Transform, transform.TransformData{
		LocalPosition: position,
		LocalRotation: originalRotation,
	})
	donburi.SetValue(airplane, component.Sprite, component.SpriteData{
		Image:            AirplaneImageByFaction(faction, 0),
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	donburi.SetValue(airplane, component.PlayerSelect, component.PlayerSelectData{
		Index:   index,
		Faction: faction,
	})

	NewCrosshair(w, airplane)

	shadow := NewStaticShadow(w, airplane)
	transform.GetTransform(shadow).LocalPosition = math.Vec2{}
}

func NewCrosshair(w donburi.World, parent *donburi.Entry) {
	crosshair := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
		),
	)

	donburi.SetValue(crosshair, transform.Transform, transform.TransformData{
		LocalScale: math.Vec2{X: 2.5, Y: 2.5},
	})

	transform.AppendChild(parent, crosshair, false)

	donburi.SetValue(crosshair, component.Sprite, component.SpriteData{
		Image:  assets.Crosshair,
		Layer:  component.SpriteLayerGroundUnits,
		Pivot:  component.SpritePivotCenter,
		Hidden: true,
	})

	label := w.Entry(
		w.Create(
			transform.Transform,
			component.Label,
		),
	)

	donburi.SetValue(label, component.Label, component.LabelData{
		Text:   "",
		Hidden: true,
	})

	transform.AppendChild(crosshair, label, false)
	transform.GetTransform(label).LocalPosition = math.Vec2{X: -25, Y: 30}
}
