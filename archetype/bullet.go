package archetype

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

const (
	playerBulletSpeed = 10
	enemyBulletSpeed  = 4
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
		),
	)

	image := assets.LaserSingle

	originalRotation := -90.0

	t := transform.Transform.Get(bullet)
	t.LocalPosition = position
	t.LocalRotation = originalRotation + localRotation

	donburi.SetValue(bullet, component.Velocity, component.VelocityData{
		Velocity: transform.Right(bullet).MulScalar(playerBulletSpeed),
	})

	donburi.SetValue(bullet, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: originalRotation,
	})

	width, height := image.Size()

	donburi.SetValue(bullet, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerPlayerBullets,
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

	image := assets.Rocket

	t := transform.Transform.Get(bullet)
	t.LocalPosition = position
	t.LocalRotation = rotation

	donburi.SetValue(bullet, component.Velocity, component.VelocityData{
		Velocity: transform.Right(bullet).MulScalar(enemyBulletSpeed),
	})

	donburi.SetValue(bullet, component.Sprite, component.SpriteData{
		Image:            image,
		Layer:            component.SpriteLayerAirUnits,
		Pivot:            component.SpritePivotCenter,
		OriginalRotation: -90,
	})

	width, height := image.Size()

	donburi.SetValue(bullet, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerEnemyBullets,
	})
}
