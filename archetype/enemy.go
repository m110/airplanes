package archetype

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewEnemyAirplane(
	w donburi.World,
	position engine.Vector,
	rotation float64,
	speed float64,
	path assets.Path,
) {
	airplane := w.Entry(
		w.Create(
			component.Transform,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
			component.Wreckable,
		),
	)

	originalRotation := -90.0

	donburi.SetValue(airplane, component.Transform, component.TransformData{
		LocalPosition: position,
		LocalRotation: originalRotation + rotation,
	})

	image := assets.AirplaneGraySmall
	donburi.SetValue(airplane, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	width, height := image.Size()

	donburi.SetValue(airplane, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerAirEnemies,
	})

	if len(path.Points) > 0 {
		donburi.SetValue(airplane, component.AI, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		donburi.SetValue(airplane, component.AI, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	donburi.SetValue(airplane, component.Health, component.HealthData{
		Health:               3,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
		DamageIndicator:      newDamageIndicator(w, airplane),
	})

	NewStaticShadow(w, airplane)
}

type spritePiece struct {
	Rect   engine.Rect
	Offset engine.Vector
}

func NewEnemyAirplaneWreck(w donburi.World, position engine.Vector, rotation float64, sprite *component.SpriteData) {
	width, height := sprite.Image.Size()
	cutpointX := float64(width) * engine.RandomRange(0.2, 0.8)
	cutpointY := float64(height) * engine.RandomRange(0.2, 0.8)

	offsetX := cutpointX / 2
	offsetY := cutpointY / 2

	pieces := []spritePiece{
		{
			Rect: engine.Rect{
				X:      0,
				Y:      0,
				Width:  cutpointX,
				Height: cutpointY,
			},
			Offset: engine.Vector{
				X: -offsetX,
				Y: -offsetY,
			},
		},
		{
			Rect: engine.Rect{
				X:      cutpointX,
				Y:      0,
				Width:  float64(width) - cutpointX,
				Height: cutpointY,
			},
			Offset: engine.Vector{
				X: offsetX,
				Y: -offsetY,
			},
		},
		{
			Rect: engine.Rect{
				X:      0,
				Y:      cutpointY,
				Width:  cutpointX,
				Height: float64(height) - cutpointY,
			},
			Offset: engine.Vector{
				X: -offsetX,
				Y: offsetY,
			},
		},
		{
			Rect: engine.Rect{
				X:      cutpointX,
				Y:      cutpointY,
				Width:  float64(width) - cutpointX,
				Height: float64(height) - cutpointY,
			},
			Offset: engine.Vector{
				X: offsetX,
				Y: offsetY,
			},
		},
	}

	// TODO need to rotate here?
	baseImage := ebiten.NewImageFromImage(sprite.Image)

	for _, p := range pieces {
		img := baseImage.SubImage(p.Rect.ToImageRectangle()).(*ebiten.Image)

		wreck := w.Entry(
			w.Create(
				component.Transform,
				component.Velocity,
				component.Altitude,
				component.Sprite,
			),
		)

		// TODO This isn't accurate
		pos := position
		pos.X += p.Offset.X
		pos.Y += p.Offset.Y

		transform := component.GetTransform(wreck)
		transform.SetWorldPosition(pos)
		transform.SetWorldRotation(rotation)

		donburi.SetValue(wreck, component.Sprite, component.SpriteData{
			Image:            img,
			Layer:            component.SpriteLayerFallingWrecks,
			Pivot:            sprite.Pivot,
			OriginalRotation: sprite.OriginalRotation,
		})

		donburi.SetValue(wreck, component.Velocity, component.VelocityData{
			Velocity: component.GetTransform(wreck).Right().MulScalar(0.75),
		})

		donburi.SetValue(wreck, component.Altitude, component.AltitudeData{
			Altitude: 1.0,
			Velocity: -0.01,
			Falling:  true,
		})

		NewStaticShadow(w, wreck)
	}
}

func NewEnemyTank(
	w donburi.World,
	position engine.Vector,
	rotation float64,
	speed float64,
	path assets.Path,
) {
	tank := w.Entry(
		w.Create(
			component.Transform,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
		),
	)

	donburi.SetValue(tank, component.Transform, component.TransformData{
		LocalPosition: position,
		LocalRotation: rotation,
	})

	image := assets.TankBase
	donburi.SetValue(tank, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerGroundUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: 0,
	})

	width, height := image.Size()

	donburi.SetValue(tank, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerGroundEnemies,
	})

	if len(path.Points) > 0 {
		donburi.SetValue(tank, component.AI, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		donburi.SetValue(tank, component.AI, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	donburi.SetValue(tank, component.Health, component.HealthData{
		Health:               5,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
		DamageIndicator:      newDamageIndicator(w, tank),
	})

	gun := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
			component.Despawnable,
			component.Observer,
			component.Shooter,
		),
	)

	originalRotation := 90.0

	donburi.SetValue(gun, component.Transform, component.TransformData{
		LocalPosition: position,
		LocalRotation: originalRotation + rotation,
	})

	donburi.SetValue(gun, component.Sprite, component.SpriteData{
		Image:            assets.TankGun,
		Layer:            component.SpriteLayerGroundUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	donburi.SetValue(gun, component.Observer, component.ObserverData{
		LookFor: query.NewQuery(filter.Contains(component.PlayerAirplane)),
	})

	donburi.SetValue(gun, component.Shooter, component.ShooterData{
		Type:       component.ShooterTypeRocket,
		ShootTimer: engine.NewTimer(time.Millisecond * 2500),
	})

	component.GetTransform(tank).AppendChild(tank, gun, true)
}

func newDamageIndicator(w donburi.World, parent *donburi.Entry) *component.SpriteData {
	indicator := w.Entry(
		w.Create(
			component.Transform,
			component.Sprite,
		),
	)

	parentSprite := component.GetSprite(parent)

	image := ebiten.NewImage(parentSprite.Image.Size())
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Translate(1, 1, 1, 0)
	image.DrawImage(parentSprite.Image, op)

	donburi.SetValue(indicator, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerIndicators,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
		Hidden:           true,
	})

	component.GetTransform(parent).AppendChild(parent, indicator, false)

	return component.GetSprite(indicator)
}
