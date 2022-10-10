package archetypes

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

const bulletSpeed = 10

func NewBullet(w donburi.World, player *component.PlayerData, position engine.Vector) {
	width := float64(assets.LaserSingle.Bounds().Dy())

	if player.WeaponLevel == component.WeaponLevelSingle ||
		player.WeaponLevel == component.WeaponLevelSingleFast {
		newBullet(w, engine.Vector{
			X: position.X,
			Y: position.Y - width,
		}, 0)
	}

	if player.WeaponLevel == component.WeaponLevelDouble ||
		player.WeaponLevel == component.WeaponLevelDoubleFast ||
		player.WeaponLevel == component.WeaponLevelDiagonal ||
		player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newBullet(w, engine.Vector{
			X: position.X - width/2,
			Y: position.Y - width/2,
		}, 0)
		newBullet(w, engine.Vector{
			X: position.X + width/2,
			Y: position.Y - width/2,
		}, 0)
	}

	if player.WeaponLevel == component.WeaponLevelDiagonal ||
		player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newBullet(w, engine.Vector{
			X: position.X - width,
			Y: position.Y - width,
		}, -30)
		newBullet(w, engine.Vector{
			X: position.X + width,
			Y: position.Y - width,
		}, 30)
	}

	if player.WeaponLevel == component.WeaponLevelDoubleDiagonal {
		newBullet(w, engine.Vector{
			X: position.X - width*1.1,
			Y: position.Y,
		}, -30)
		newBullet(w, engine.Vector{
			X: position.X + width*1.1,
			Y: position.Y,
		}, 30)
	}
}

func newBullet(w donburi.World, position engine.Vector, rotation float64) {
	bullet := w.Entry(
		w.Create(
			component.Velocity,
			component.Transform,
			component.Sprite,
			component.Despawnable,
			component.Collider,
		),
	)

	image := assets.LaserSingle

	originalRotation := -90.0

	donburi.SetValue(bullet, component.Transform, component.TransformData{
		Position: position,
		Rotation: originalRotation + rotation,
	})

	donburi.SetValue(bullet, component.Velocity, component.VelocityData{
		Velocity: component.GetTransform(bullet).Right().MulScalar(bulletSpeed),
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
		Layer:  component.CollisionLayerBullets,
	})
}
