package archetype

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

const (
	playerBulletSpeed = 10
	enemyBulletSpeed  = 4
	enemyMissileSpeed = 2
)

func NewPlayerBullet(w donburi.World, player *component.PlayerData, position math.Vec2) {
	width := float64(assets.LaserSingle.Bounds().Dy())

	if player.WeaponLevel == component.WeaponLevelSingle ||
		player.WeaponLevel == component.WeaponLevelSingleFast {
		newPlayerBullet(w, math.Vec2{
			X: position.X,
			Y: position.Y - width,
		}, 0)
	}

	if player.WeaponLevel == component.WeaponLevelDouble ||
		player.WeaponLevel == component.WeaponLevelDoubleFast ||
		player.WeaponLevel == component.WeaponLevelDiagonal ||
		player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newPlayerBullet(w, math.Vec2{
			X: position.X - width/2,
			Y: position.Y - width/2,
		}, 0)
		newPlayerBullet(w, math.Vec2{
			X: position.X + width/2,
			Y: position.Y - width/2,
		}, 0)
	}

	if player.WeaponLevel == component.WeaponLevelDiagonal ||
		player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newPlayerBullet(w, math.Vec2{
			X: position.X - width,
			Y: position.Y - width,
		}, -30)
		newPlayerBullet(w, math.Vec2{
			X: position.X + width,
			Y: position.Y - width,
		}, 30)
	}

	if player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newPlayerBullet(w, math.Vec2{
			X: position.X - width*1.1,
			Y: position.Y,
		}, -30)
		newPlayerBullet(w, math.Vec2{
			X: position.X + width*1.1,
			Y: position.Y,
		}, 30)
	}
}

func newPlayerBullet(w donburi.World, position math.Vec2, localRotation float64) {
	bullet := w.Entry(
		w.Create(
			component.Velocity,
			transform.Transform,
			component.Sprite,
			component.Despawnable,
			component.Collider,
			component.DistanceLimit,
		),
	)

	image := assets.LaserSingle

	originalRotation := -90.0

	t := transform.Transform.Get(bullet)
	t.LocalPosition = position
	t.LocalRotation = originalRotation + localRotation

	component.Velocity.SetValue(bullet, component.VelocityData{
		Velocity: transform.Right(bullet).MulScalar(playerBulletSpeed),
	})

	component.Sprite.SetValue(bullet, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	width, height := image.Size()

	component.Collider.SetValue(bullet, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerPlayerBullets,
	})

	component.DistanceLimit.SetValue(bullet, component.DistanceLimitData{
		MaxDistance: 200,
	})
}

func NewEnemyBullet(w donburi.World, position math.Vec2, rotation float64) {
	bullet := w.Entry(
		w.Create(
			component.Velocity,
			transform.Transform,
			component.Sprite,
			component.Despawnable,
			component.Collider,
		),
	)

	image := assets.Bullet

	t := transform.Transform.Get(bullet)
	t.LocalPosition = position
	t.LocalRotation = rotation

	component.Velocity.SetValue(bullet, component.VelocityData{
		Velocity: transform.Right(bullet).MulScalar(enemyBulletSpeed),
	})

	component.Sprite.SetValue(bullet, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: -90,
	})

	width, height := image.Size()

	component.Collider.SetValue(bullet, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerEnemyBullets,
	})
}

func NewEnemyMissile(w donburi.World, position math.Vec2, rotation float64) {
	missile := w.Entry(
		w.Create(
			component.Velocity,
			transform.Transform,
			component.Sprite,
			component.Despawnable,
			component.Collider,
			component.Follower,
		),
	)

	image := assets.Missile

	t := transform.Transform.Get(missile)
	t.LocalPosition = position
	t.LocalRotation = rotation

	component.Velocity.SetValue(missile, component.VelocityData{
		Velocity: transform.Right(missile).MulScalar(enemyMissileSpeed),
	})

	component.Sprite.SetValue(missile, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: -90,
	})

	width, height := image.Size()

	component.Collider.SetValue(missile, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerEnemyBullets,
	})

	component.Follower.SetValue(missile, component.FollowerData{
		Target:         component.ClosestTarget(w, missile, donburi.NewQuery(filter.Contains(component.PlayerAirplane))),
		FollowingSpeed: enemyMissileSpeed,
		FollowingTimer: engine.NewTimer(3 * time.Second),
	})
}
