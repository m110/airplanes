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
) *donburi.Entry {
	enemy := w.Entry(
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

	donburi.SetValue(enemy, component.Position, component.PositionData{
		Position: position,
	})
	donburi.SetValue(enemy, component.Rotation, component.RotationData{
		Angle:         rotation,
		OriginalAngle: -90,
	})

	image := assets.AirplaneGraySmall
	donburi.SetValue(enemy, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerAirUnits,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()

	donburi.SetValue(enemy, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerAirEnemies,
	})

	if len(path.Points) > 0 {
		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	donburi.SetValue(enemy, component.Health, component.HealthData{
		Health:               3,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
	})

	return enemy
}

func NewEnemyTank(
	w donburi.World,
	position engine.Vector,
	rotation float64,
	speed float64,
	path assets.Path,
) *donburi.Entry {
	enemy := w.Entry(
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

	donburi.SetValue(enemy, component.Position, component.PositionData{
		Position: position,
	})
	donburi.SetValue(enemy, component.Rotation, component.RotationData{
		Angle:         rotation,
		OriginalAngle: 0,
	})

	image := assets.TankBase
	donburi.SetValue(enemy, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerGroundUnits,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()

	donburi.SetValue(enemy, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerGroundEnemies,
	})

	if len(path.Points) > 0 {
		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	donburi.SetValue(enemy, component.Health, component.HealthData{
		Health:               5,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
	})

	return enemy
}
