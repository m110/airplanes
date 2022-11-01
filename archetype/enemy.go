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

	transform.GetTransform(spawn).LocalPosition = position
	component.GetSpawnable(spawn).SpawnFunc = spawnFunc
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

	donburi.SetValue(airplane, transform.Transform, transform.TransformData{
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
	t := transform.GetTransform(tank)
	t.LocalPosition = position
	t.LocalRotation = rotation

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
			transform.Transform,
			component.Sprite,
			component.Despawnable,
			component.Observer,
			component.Shooter,
		),
	)

	originalRotation := 90.0
	transform.Reset(gun)
	gunT := transform.GetTransform(gun)
	gunT.LocalPosition = position
	gunT.LocalRotation = originalRotation + rotation

	donburi.SetValue(gun, component.Sprite, component.SpriteData{
		Image:            assets.TankGun,
		Layer:            component.SpriteLayerGroundGuns,
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

	transform.AppendChild(tank, gun, true)
}

func newDamageIndicator(w donburi.World, parent *donburi.Entry) *component.SpriteData {
	indicator := w.Entry(
		w.Create(
			transform.Transform,
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

	transform.AppendChild(parent, indicator, false)

	return component.GetSprite(indicator)
}
