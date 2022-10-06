package system

import (
	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type AI struct {
	query *query.Query
}

func NewAI() *AI {
	return &AI{
		query: query.NewQuery(
			filter.Contains(
				component.Position,
				component.Velocity,
				component.Sprite,
				component.AI,
			),
		),
	}
}

func (a *AI) Update(w donburi.World) {
	a.query.EachEntity(w, func(entry *donburi.Entry) {
		ai := component.GetAI(entry)

		if !ai.Spawned {
			cameraPos := component.GetPosition(archetypes.MustFindCamera(w))
			position := component.GetPosition(entry)
			sprite := component.GetSprite(entry)

			if position.Y+float64(sprite.Image.Bounds().Dy()) > cameraPos.Y {
				ai.Spawned = true

				velocity := component.GetVelocity(entry)

				switch ai.Type {
				case component.AITypeConstantVelocity:
					velocity.Y = ai.ConstantVelocity
				case component.AITypeFollowPath:
					// TODO
				}
			}
		}
	})
}
