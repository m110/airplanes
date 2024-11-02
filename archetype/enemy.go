package archetype

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

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

	health := component.Health.Get(airplane)
	health.Health = 3
	health.DamageIndicator = newDamageIndicator(w, airplane)

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

	health := component.Health.Get(tank)
	health.Health = 5
	health.DamageIndicator = newDamageIndicator(w, tank)

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
		LookFor: donburi.NewQuery(filter.Contains(component.PlayerAirplane)),
	})

	component.Shooter.SetValue(gun, component.ShooterData{
		Type:       component.ShooterTypeBullet,
		ShootTimer: engine.NewTimer(time.Millisecond * 2500),
	})

	transform.AppendChild(tank, gun, true)
}

func NewEnemyTurretBeam(
	w donburi.World,
	position math.Vec2,
	rotation float64,
) {
	turret := newEnemyTurret(w, position, rotation)

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
	gunT := transform.Transform.Get(gun)
	gunT.LocalPosition = position
	gunT.LocalRotation = originalRotation + rotation

	component.Sprite.SetValue(gun, component.SpriteData{
		Image:            assets.TurretGunSingle,
		Layer:            component.SpriteLayerGroundGuns,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	component.Observer.SetValue(gun, component.ObserverData{
		LookFor: donburi.NewQuery(filter.Contains(component.PlayerAirplane)),
	})

	component.Shooter.SetValue(gun, component.ShooterData{
		Type:       component.ShooterTypeBeam,
		ShootTimer: engine.NewTimer(time.Millisecond * 5000),
	})

	transform.AppendChild(turret, gun, true)
}

func NewEnemyTurretMissiles(
	w donburi.World,
	position math.Vec2,
	rotation float64,
) {
	turret := newEnemyTurret(w, position, rotation)

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
	gunT := transform.Transform.Get(gun)
	gunT.LocalPosition = position
	gunT.LocalRotation = originalRotation + rotation

	component.Sprite.SetValue(gun, component.SpriteData{
		Image:            assets.TurretGunDouble,
		Layer:            component.SpriteLayerGroundGuns,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	component.Observer.SetValue(gun, component.ObserverData{
		LookFor: donburi.NewQuery(filter.Contains(component.PlayerAirplane)),
	})

	component.Shooter.SetValue(gun, component.ShooterData{
		Type:       component.ShooterTypeMissile,
		ShootTimer: engine.NewTimer(time.Millisecond * 5000),
	})

	transform.AppendChild(turret, gun, true)
}

func newEnemyTurret(
	w donburi.World,
	position math.Vec2,
	rotation float64,
) *donburi.Entry {
	turret := w.Entry(
		w.Create(
			transform.Transform,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
			component.Health,
		),
	)

	t := transform.Transform.Get(turret)
	t.LocalPosition = position
	t.LocalRotation = rotation

	image := assets.TurretBase
	component.Sprite.SetValue(turret, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerGroundUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: 0,
	})

	width, height := image.Size()

	component.Collider.SetValue(turret, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerGroundEnemies,
	})

	health := component.Health.Get(turret)
	health.Health = 5
	health.DamageIndicator = newDamageIndicator(w, turret)

	return turret
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
