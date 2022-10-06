package archetypes

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewEnemy(w donburi.World, position component.PositionData, rotation float64) *donburi.Entry {
	enemy := w.Entry(
		w.Create(
			component.Position,
			component.Rotation,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
		),
	)

	donburi.SetValue(enemy, component.Position, position)
	component.GetRotation(enemy).Angle = rotation
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
