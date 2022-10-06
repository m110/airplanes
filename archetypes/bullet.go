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
		),
	)

	component.GetVelocity(bullet).Y = -10

	sprite := component.GetSprite(bullet)
	sprite.Image = assets.LaserSingle
	sprite.Layer = component.SpriteLayerUnits

	component.GetDespawnable(bullet).Spawned = true

	return bullet
}
