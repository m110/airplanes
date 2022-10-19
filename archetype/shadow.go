package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

const (
	shadowColorScale = 0.5
	shadowColorAlpha = 0.4
)

func NewDynamicShadow(w donburi.World, parent *donburi.Entry) *donburi.Entry {
	return newShadow(w, parent, true)
}

func NewStaticShadow(w donburi.World, parent *donburi.Entry) *donburi.Entry {
	return newShadow(w, parent, false)
}

func newShadow(w donburi.World, parent *donburi.Entry, dynamic bool) *donburi.Entry {
	shadow := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
			component.ShadowTag,
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

	var img *ebiten.Image
	if dynamic {
		img = ebiten.NewImageFromImage(parentSprite.Image)
	} else {
		img = ebiten.NewImage(parentSprite.Image.Size())
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(0, 0, 0, shadowColorAlpha)
		op.ColorM.Translate(shadowColorScale, shadowColorScale, shadowColorScale, 0)
		img.DrawImage(parentSprite.Image, op)
	}

	spriteData := component.SpriteData{
		Image:            img,
		Layer:            component.SpriteLayerShadows,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
	}

	if dynamic {
		spriteData.ColorOverride = &component.ColorOverride{
			R: shadowColorScale,
			G: shadowColorScale,
			B: shadowColorScale,
			A: shadowColorAlpha,
		}
	}

	donburi.SetValue(shadow, component.Sprite, spriteData)

	return shadow
}
