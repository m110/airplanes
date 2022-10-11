package system

import (
	"math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type AI struct {
	query *query.Query
}

func NewAI() *AI {
	return &AI{
		query: query.NewQuery(
			filter.Contains(
				component.Transform,
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
		if ai.Spawned {
			if ai.Type == component.AITypeFollowPath {
				if ai.NextTarget >= len(ai.Path) {
					return
				}

				position := component.GetTransform(entry).LocalPosition
				velocity := component.GetVelocity(entry)

				target := ai.Path[ai.NextTarget]

				x := target.X - position.X
				y := target.Y - position.Y

				dist := math.Sqrt(x*x + y*y)
				if dist < 1 {
					ai.NextTarget++
					if ai.PathLoops && ai.NextTarget >= len(ai.Path) {
						ai.NextTarget = 0
					}
					return
				}

				transform := component.GetTransform(entry)

				// TODO Could be simplified perhaps ^^'
				angle := math.Round(engine.ToDegrees(math.Atan2(y, x)))

				// TODO Learn trigonometry
				if transform.WorldRotation()-angle > 180.0 {
					angle = float64(int(angle+360.0) % 360)
				} else if transform.WorldRotation()-angle < -180.0 {
					angle = float64(int(angle-360.0) % 360)
				}

				maxRotation := 2.0 * ai.Speed
				targetAngle := angle

				diff := targetAngle - transform.WorldRotation()
				if math.Abs(diff) > maxRotation {
					if diff > 0 {
						diff = maxRotation
					} else {
						diff = -maxRotation
					}
				}
				transform.LocalRotation += diff

				// TODO Should use transform.Right() instead but it doesn't work
				radians := engine.ToRadians(angle)
				velocity.Velocity.X = math.Cos(radians) * ai.Speed
				velocity.Velocity.Y = math.Sin(radians) * ai.Speed
			}
		} else {
			spawnEnemy(w, entry)
		}
	})
}

func spawnEnemy(w donburi.World, entry *donburi.Entry) {
	cameraPos := component.GetTransform(archetypes.MustFindCamera(w)).WorldPosition()

	ai := component.GetAI(entry)
	transform := component.GetTransform(entry)
	sprite := component.GetSprite(entry)

	if transform.LocalPosition.Y+float64(sprite.Image.Bounds().Dy()) > cameraPos.Y {
		ai.Spawned = true

		velocity := component.GetVelocity(entry)

		if ai.Type == component.AITypeConstantVelocity {
			velocity.Velocity = transform.Right().MulScalar(ai.Speed)
		}
	}
}
