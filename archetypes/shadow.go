package archetypes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewShadow(w donburi.World, parent *donburi.Entry) {
	shadow := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
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

	shadowImage := ebiten.NewImage(parentSprite.Image.Size())
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(0, 0, 0, 0.4)
	op.ColorM.Translate(0.5, 0.5, 0.5, 0)
	shadowImage.DrawImage(parentSprite.Image, op)

	donburi.SetValue(shadow, component.Sprite, component.SpriteData{
		Image:            shadowImage,
		Layer:            component.SpriteLayerShadows,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
	})
}
