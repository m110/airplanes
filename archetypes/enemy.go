package archetypes

import (
	"time"

	"github.com/yohamta/donburi"

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
			component.Position,
			component.Rotation,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
		),
	)

	donburi.SetValue(airplane, component.Position, component.PositionData{
		Position: position,
	})
	donburi.SetValue(airplane, component.Rotation, component.RotationData{
		Angle:         rotation,
		OriginalAngle: -90,
	})

	image := assets.AirplaneGraySmall
	donburi.SetValue(airplane, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerAirUnits,
		Pivot: component.SpritePivotCenter,
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
	})
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
			component.Position,
			component.Rotation,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
		),
	)

	donburi.SetValue(tank, component.Position, component.PositionData{
		Position: position,
	})
	donburi.SetValue(tank, component.Rotation, component.RotationData{
		Angle:         rotation,
		OriginalAngle: 0,
	})

	image := assets.TankBase
	donburi.SetValue(tank, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerGroundUnits,
		Pivot: component.SpritePivotCenter,
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
	})

	gun := w.Entry(
		w.Create(
			component.Position,
			component.Rotation,
			component.Sprite,
			component.Despawnable,
		),
	)

	donburi.SetValue(gun, component.Position, component.PositionData{
		Position: position,
	})

	donburi.SetValue(gun, component.Rotation, component.RotationData{
		Angle:         rotation,
		OriginalAngle: 90,
	})

	donburi.SetValue(gun, component.Sprite, component.SpriteData{
		Image: assets.TankGun,
		Layer: component.SpriteLayerGroundUnits,
		Pivot: component.SpritePivotCenter,
	})

	component.GetPosition(tank).AppendChild(tank, gun)
}
