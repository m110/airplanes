package archetypes

import (
	"time"

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
	})

	NewShadow(w, airplane)
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

	component.GetTransform(tank).AppendChild(tank, gun)
}
