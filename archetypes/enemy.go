package archetypes

import (
	"github.com/samber/lo"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewEnemy(
	w donburi.World,
	position component.PositionData,
	rotation float64,
	speed float64,
	path []assets.Position,
) *donburi.Entry {
	enemy := w.Entry(
		w.Create(
			component.Position,
			component.Rotation,
			component.Velocity,
			component.Sprite,
			component.AI,
			component.Despawnable,
			component.Collider,
		),
	)

	donburi.SetValue(enemy, component.Position, position)
	component.GetRotation(enemy).Angle = rotation

	image := assets.ShipGraySmall
	donburi.SetValue(enemy, component.Sprite, component.SpriteData{
		Image: image,
		Layer: component.SpriteLayerUnits,
	})

	width, height := image.Size()

	donburi.SetValue(enemy, component.Collider, component.ColliderData{
		Width:  width,
		Height: height,
		Layer:  component.CollisionLayerEnemies,
	})

	if len(path) > 0 {
		componentPath := lo.Map(path, func(p assets.Position, _ int) component.PathPosition {
			return component.PathPosition(p)
		})

		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:  component.AITypeFollowPath,
			Speed: speed,
			Path:  componentPath,
		})
	} else {
		donburi.SetValue(enemy, component.AI, component.AIData{
			Type:  component.AITypeConstantVelocity,
			Speed: speed,
		})
	}

	return enemy
}
