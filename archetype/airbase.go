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

	t := transform.Transform.Get(airplane)
	t.LocalPosition = position
	t.LocalRotation = originalRotation

	component.Sprite.SetValue(airplane, component.SpriteData{
		Image:            AirplaneImageByFaction(faction, 0),
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	component.PlayerSelect.SetValue(airplane, component.PlayerSelectData{
		Index:   index,
		Faction: faction,
	})

	NewCrosshair(w, airplane)

	shadow := NewShadow(w, airplane)
	transform.Transform.Get(shadow).LocalPosition = math.Vec2{}
}

func NewCrosshair(w donburi.World, parent *donburi.Entry) {
	crosshair := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
		),
	)

	transform.Transform.Get(crosshair).LocalScale = math.Vec2{X: 2.5, Y: 2.5}

	transform.AppendChild(parent, crosshair, false)

	component.Sprite.SetValue(crosshair, component.SpriteData{
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

	component.Label.SetValue(label, component.LabelData{
		Text:   "",
		Hidden: true,
	})

	transform.AppendChild(crosshair, label, false)
	transform.Transform.Get(label).LocalPosition = math.Vec2{X: -25, Y: 30}
}
