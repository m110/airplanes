package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

var ShadowTag = donburi.NewTag()

func NewShadow(w donburi.World, parent *donburi.Entry) *donburi.Entry {
	shadow := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
			ShadowTag,
		),
	)

	parentTransform := component.GetTransform(parent)
	parentTransform.AppendChild(parent, shadow, false)

	parentSprite := component.GetSprite(parent)
	width, height := parentSprite.Image.Size()

	transform := component.GetTransform(shadow)
	transform.LocalPosition = engine.Vector{
		X: -float64(width) * 0.35,
		Y: float64(height) * 0.35,
	}

	donburi.SetValue(shadow, component.Sprite, component.SpriteData{
		Image:            ebiten.NewImageFromImage(parentSprite.Image),
		Layer:            component.SpriteLayerShadows,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
		// TODO useful for dynamic shadows but for static ones the color transformation could be done just once
		ColorOverride: &component.ColorOverride{
			R: 0.5,
			G: 0.5,
			B: 0.5,
			A: 0.4,
		},
	})

	return shadow
}
