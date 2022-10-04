package archetypes

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewBullet(w donburi.World) *donburi.Entry {
	bulletEntity := w.Create(
		component.Velocity,
		component.Position,
		component.Sprite,
	)
	bullet := w.Entry(bulletEntity)

	donburi.SetValue(bullet, component.Position, component.PositionData{})
	donburi.SetValue(bullet, component.Velocity, component.VelocityData{
		Y: -10,
	})
	donburi.SetValue(bullet, component.Sprite, component.SpriteData{
		Image: assets.LaserSingle,
	})

	return bullet
}
