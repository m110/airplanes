package archetypes

import (
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/yohamta/donburi"
)

func NewEnemy(w donburi.World, position component.PositionData) *donburi.Entry {
	enemy := w.Entry(
		w.Create(
			component.Position,
			component.Velocity,
			component.Sprite,
			component.AI,
		),
	)

	donburi.SetValue(enemy, component.Position, position)
	donburi.SetValue(enemy, component.Sprite, component.SpriteData{
		Image: assets.ShipGraySmall,
		Layer: component.SpriteLayerUnits,
	})
	donburi.SetValue(enemy, component.AI, component.AIData{
		Type:             component.AITypeConstantVelocity,
		ConstantVelocity: 1,
	})

	return enemy
}
