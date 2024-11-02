package archetype

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewAirplaneWreck(w donburi.World, parent *donburi.Entry, sprite *component.SpriteData) {
	widthInt, heightInt := sprite.Image.Size()
	width, height := float64(widthInt), float64(heightInt)
	cutpointX := float64(int(width * engine.RandomFloatRange(0.3, 0.7)))
	cutpointY := float64(int(height * engine.RandomFloatRange(0.3, 0.7)))

	pieces := []engine.Rect{
		{
			X:      0,
			Y:      0,
			Width:  cutpointX,
			Height: cutpointY,
		},
		{
			X:      cutpointX,
			Y:      0,
			Width:  width - cutpointX,
			Height: cutpointY,
		},
		{
			X:      0,
			Y:      cutpointY,
			Width:  cutpointX,
			Height: height - cutpointY,
		},
		{
			X:      cutpointX,
			Y:      cutpointY,
			Width:  width - cutpointX,
			Height: height - cutpointY,
		},
	}

	halfW := width / 2
	halfH := height / 2

	// Rotate the base image
	baseImage := ebiten.NewImage(sprite.Image.Size())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(float64(int(transform.WorldRotation(parent)-sprite.OriginalRotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(halfW, halfH)
	baseImage.DrawImage(sprite.Image, op)

	basePos := transform.WorldPosition(parent)
	if sprite.Pivot == component.SpritePivotCenter {
		basePos.X -= halfW
		basePos.Y -= halfH
	}

	for _, p := range pieces {
		img := baseImage.SubImage(p.ToImageRectangle()).(*ebiten.Image)

		wreck := w.Entry(
			w.Create(
				transform.Transform,
				component.Velocity,
				component.Altitude,
				component.Sprite,
			),
		)

		pos := basePos
		pos.X += p.X + p.Width/2
		pos.Y += p.Y + p.Height/2

		transform.Transform.Get(wreck).LocalPosition = pos

		component.Sprite.SetValue(wreck, component.SpriteData{
			Image: img,
			Layer: component.SpriteLayerFallingWrecks,
			Pivot: sprite.Pivot,
		})

		velocity := transform.Right(parent)
		velocity.X *= engine.RandomFloatRange(0.5, 0.8)
		velocity.Y *= engine.RandomFloatRange(0.5, 0.8)

		component.Velocity.SetValue(wreck, component.VelocityData{
			Velocity:         velocity,
			RotationVelocity: engine.RandomFloatRange(-2, 2),
		})

		component.Altitude.SetValue(wreck, component.AltitudeData{
			Altitude: 1.0,
			Velocity: -0.01,
			Falling:  true,
		})

		NewShadow(w, wreck)
	}
}
