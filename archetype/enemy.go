package archetype

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewEnemySpawn(w donburi.World, position math.Vec2, spawnFunc component.SpawnFunc) {
	spawn := w.Entry(
		w.Create(
			transform.Transform,
			component.Spawnable,
		),
	)

	transform.Transform.Get(spawn).LocalPosition = position
	component.Spawnable.Get(spawn).SpawnFunc = spawnFunc
}

func NewEnemyAirplane(
	w donburi.World,
	position math.Vec2,
	rotation float64,
	speed float64,
	path assets.Path,
) {
	airplane := w.Entry(
		w.Create(
			transform.Transform,
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

	t := transform.Transform.Get(airplane)
	t.LocalPosition = position
	t.LocalRotation = originalRotation + rotation

	image := assets.AirplaneGraySmall
	component.Sprite.SetValue(airplane, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	width, height := image.Size()

	component.Collider.SetValue(airplane, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerAirEnemies,
	})

	if len(path.Points) > 0 {
		component.AI.SetValue(airplane, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		component.AI.SetValue(airplane, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	component.Health.SetValue(airplane, component.HealthData{
		Health:               3,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
		DamageIndicator:      newDamageIndicator(w, airplane),
	})

	NewShadow(w, airplane)
}

func NewEnemyTank(
	w donburi.World,
	position math.Vec2,
	rotation float64,
	speed float64,
	path assets.Path,
) {
	tank := w.Entry(
		w.Create(
			transform.Transform,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
		),
	)

	transform.Reset(tank)
	t := transform.Transform.Get(tank)
	t.LocalPosition = position
	t.LocalRotation = rotation

	image := assets.TankBase
	component.Sprite.SetValue(tank, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerGroundUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: 0,
	})

	width, height := image.Size()

	component.Collider.SetValue(tank, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerGroundEnemies,
	})

	if len(path.Points) > 0 {
		component.AI.SetValue(tank, component.AIData{
			Type:      component.AITypeFollowPath,
			Speed:     speed,
			Path:      path.Points,
			PathLoops: path.Loops,
		})
	} else {
		component.AI.SetValue(tank, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	component.Health.SetValue(tank, component.HealthData{
		Health:               5,
		DamageIndicatorTimer: engine.NewTimer(time.Millisecond * 100),
		DamageIndicator:      newDamageIndicator(w, tank),
	})

	gun := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
			component.Despawnable,
			component.Observer,
			component.Shooter,
		),
	)

	originalRotation := 90.0
	transform.Reset(gun)
	gunT := transform.Transform.Get(gun)
	gunT.LocalPosition = position
	gunT.LocalRotation = originalRotation + rotation

	component.Sprite.SetValue(gun, component.SpriteData{
		Image:            assets.TankGun,
		Layer:            component.SpriteLayerGroundGuns,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	component.Observer.SetValue(gun, component.ObserverData{
		LookFor: query.NewQuery(filter.Contains(component.PlayerAirplane)),
	})

	component.Shooter.SetValue(gun, component.ShooterData{
		Type:       component.ShooterTypeRocket,
		ShootTimer: engine.NewTimer(time.Millisecond * 2500),
	})

	transform.AppendChild(tank, gun, true)
}

func newDamageIndicator(w donburi.World, parent *donburi.Entry) *component.SpriteData {
	indicator := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
		),
	)

	parentSprite := component.Sprite.Get(parent)

	image := ebiten.NewImage(parentSprite.Image.Size())
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Translate(1, 1, 1, 0)
	image.DrawImage(parentSprite.Image, op)

	component.Sprite.SetValue(indicator, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerIndicators,
		Pivot:            parentSprite.Pivot,
		OriginalRotation: parentSprite.OriginalRotation,
		Hidden:           true,
	})

	transform.AppendChild(parent, indicator, false)

	return component.Sprite.Get(indicator)
}
