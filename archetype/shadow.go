package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/component"
)

const (
	shadowColorScale = 0.5
	shadowColorAlpha = 0.4

	// TODO Should this be based on the sprite's width?
	MaxShadowPosition = 12
)

func NewShadow(w donburi.World, parent *donburi.Entry) *donburi.Entry {
	shadow := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
			component.ShadowTag,
		),
	)

	transform.AppendChild(parent, shadow, false)

	transform := transform.Transform.Get(shadow)
	transform.LocalPosition = math.Vec2{
		X: -MaxShadowPosition,
		Y: MaxShadowPosition,
	}

	parentSprite := component.Sprite.Get(parent)

	spriteData := component.SpriteData{
		Image:            ShadowImage(parentSprite.Image),
		Layer:            component.SpriteLayerShadows,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
	}

	component.Sprite.SetValue(shadow, spriteData)

	return shadow
}

func ShadowImage(source *ebiten.Image) *ebiten.Image {
	shadow := ebiten.NewImage(source.Size())
	op := &ebiten.DrawImageOptions{}
	ShadowDrawOptions(op)
	shadow.DrawImage(source, op)
	return shadow
}

func ShadowDrawOptions(op *ebiten.DrawImageOptions) {
	op.ColorM.Scale(0, 0, 0, shadowColorAlpha)
	op.ColorM.Translate(shadowColorScale, shadowColorScale, shadowColorScale, 0)
}
