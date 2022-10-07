package system

import (
	"math"

	"github.com/m110/airplanes/archetypes"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
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
				component.Rotation,
			),
		),
	}
}

func (a *AI) Update(w donburi.World) {
	a.query.EachEntity(w, func(entry *donburi.Entry) {
		ai := component.GetAI(entry)
		if ai.Spawned {
			if ai.Type == component.AITypeFollowPath {
				if ai.NextTarget >= len(ai.Path) {
					return
				}

				position := component.GetPosition(entry)
				velocity := component.GetVelocity(entry)

				target := ai.Path[ai.NextTarget]

				x := target.X - position.X
				y := target.Y - position.Y

				dist := math.Sqrt(x*x + y*y)
				if dist < 1 {
					ai.NextTarget++
					return
				}

				// TODO Could be simplified perhaps ^^'
				angle := math.Atan2(y, x) * 180.0 / math.Pi
				rotation := component.GetRotation(entry)

				maxRotation := 2.0 * ai.Speed
				targetAngle := angle + 90
				diff := targetAngle - rotation.Angle
				if math.Abs(diff) > maxRotation {
					if diff > 0 {
						diff = maxRotation
					} else {
						diff = -maxRotation
					}
				}
				rotation.Angle += diff

				radians := float64(angle) / 180.0 * math.Pi

				velocity.X = math.Cos(radians) * ai.Speed
				velocity.Y = math.Sin(radians) * ai.Speed
			}
		} else {
			spawnEnemy(w, entry)
		}
	})
}

func spawnEnemy(w donburi.World, entry *donburi.Entry) {
	cameraPos := component.GetPosition(archetypes.MustFindCamera(w))

	ai := component.GetAI(entry)
	position := component.GetPosition(entry)
	rotation := component.GetRotation(entry)
	sprite := component.GetSprite(entry)

	if position.Y+float64(sprite.Image.Bounds().Dy()) > cameraPos.Y {
		ai.Spawned = true

		velocity := component.GetVelocity(entry)

		if ai.Type == component.AITypeConstantVelocity {
			radians := float64(rotation.Angle-90) / 180.0 * math.Pi
			velocity.X = math.Cos(radians) * ai.Speed
			velocity.Y = math.Sin(radians) * ai.Speed
		}
	}
}
