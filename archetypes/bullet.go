package archetypes

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewBullet(w donburi.World) *donburi.Entry {
	bullet := w.Entry(
		w.Create(
			component.Velocity,
			component.Position,
			component.Sprite,
			component.Despawnable,
			component.Collider,
		),
	)

	component.GetVelocity(bullet).Y = -10

	image := assets.LaserSingle

	donburi.SetValue(bullet, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerUnits,
		Pivot: component.SpritePivotCenter,
	})

	width, height := image.Size()

	donburi.SetValue(bullet, component.Collider, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerBullets,
	})

	return bullet
}
