package system

import (
	"math"

	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type AI struct {
	query *donburi.Query
}

func NewAI() *AI {
	return &AI{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Velocity,
				component.AI,
			),
		),
	}
}

func (a *AI) Update(w donburi.World) {
	a.query.Each(w, func(entry *donburi.Entry) {
		ai := component.AI.Get(entry)
		if ai.Type == component.AITypeFollowPath {
			if ai.NextTarget >= len(ai.Path) {
				return
			}

			position := transform.Transform.Get(entry).LocalPosition
			velocity := component.Velocity.Get(entry)

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

			transform.Transform.Get(entry).LocalRotation += diff

			// TODO Should use transform.Right() instead but it doesn't work
			radians := dmath.ToRadians(angle)
			velocity.Velocity.X = math.Cos(radians) * ai.Speed
			velocity.Velocity.Y = math.Sin(radians) * ai.Speed
		} else if ai.Type == component.AITypeConstantVelocity && !ai.StartedMoving {
			velocity := component.Velocity.Get(entry)
			velocity.Velocity = transform.Right(entry).MulScalar(ai.Speed)
			ai.StartedMoving = true
		}
	})
}
