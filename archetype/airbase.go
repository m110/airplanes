package archetype

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewAirbaseAirplane(w donburi.World, position engine.Vector, faction component.PlayerFaction, index int) {
	airplane := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
			component.Velocity,
			component.PlayerSelect,
			component.Altitude,
		),
	)

	originalRotation := -90.0

	donburi.SetValue(airplane, component.Transform, component.TransformData{
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
	component.GetTransform(shadow).LocalPosition = engine.Vector{}
}

func NewCrosshair(w donburi.World, parent *donburi.Entry) {
	crosshair := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
		),
	)

	donburi.SetValue(crosshair, component.Transform, component.TransformData{
		LocalScale: engine.Vector{X: 2.5, Y: 2.5},
	})

	component.GetTransform(parent).AppendChild(parent, crosshair, false)

	donburi.SetValue(crosshair, component.Sprite, component.SpriteData{
		Image:  assets.Crosshair,
		Layer:  component.SpriteLayerGroundUnits,
		Pivot:  component.SpritePivotCenter,
		Hidden: true,
	})

	label := w.Entry(
		w.Create(
			component.Transform,
			component.Label,
		),
	)

	donburi.SetValue(label, component.Label, component.LabelData{
		Text:   "",
		Hidden: true,
	})

	component.GetTransform(crosshair).AppendChild(crosshair, label, false)
	component.GetTransform(label).LocalPosition = engine.Vector{X: -25, Y: 30}
}
