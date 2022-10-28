package system

import (
	"math"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type AI struct {
	query *query.Query
}

func NewAI() *AI {
	return &AI{
		query: query.NewQuery(
			filter.Contains(
				transform.Transform,
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

				position := transform.GetTransform(entry).LocalPosition
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

				worldRotation := transform.WorldRotation(entry)

				// TODO Could be simplified perhaps ^^'
				angle := math.Round(dmath.ToDegrees(math.Atan2(y, x)))

				// TODO Learn trigonometry
				if worldRotation-angle > 180.0 {
					angle = float64(int(angle+360.0) % 360)
				} else if worldRotation-angle < -180.0 {
					angle = float64(int(angle-360.0) % 360)
				}

				maxRotation := 2.0 * ai.Speed
				targetAngle := angle

				diff := targetAngle - worldRotation
				if math.Abs(diff) > maxRotation {
					if diff > 0 {
						diff = maxRotation
					} else {
						diff = -maxRotation
					}
				}

				transform.GetTransform(entry).LocalRotation += diff

				// TODO Should use transform.Right() instead but it doesn't work
				radians := dmath.ToRadians(angle)
				velocity.Velocity.X = math.Cos(radians) * ai.Speed
				velocity.Velocity.Y = math.Sin(radians) * ai.Speed
			}
		} else {
			spawnEnemy(w, entry)
		}
	})
}

func spawnEnemy(w donburi.World, entry *donburi.Entry) {
	cameraPos := transform.WorldPosition(archetype.MustFindCamera(w))

	ai := component.GetAI(entry)
	sprite := component.GetSprite(entry)
	t := transform.GetTransform(entry)

	if t.LocalPosition.Y+float64(sprite.Image.Bounds().Dy()) > cameraPos.Y {
		ai.Spawned = true

		velocity := component.GetVelocity(entry)

		if ai.Type == component.AITypeConstantVelocity {
			vel := transform.Right(entry)
			vel.MulScalar(ai.Speed)
			velocity.Velocity = vel
		}
	}
}
