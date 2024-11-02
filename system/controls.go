package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Controls struct {
	query *donburi.Query
}

func NewControls() *Controls {
	return &Controls{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Input,
				component.Velocity,
				component.Sprite,
				component.PlayerAirplane,
			),
		),
	}
}

func (i *Controls) Update(w donburi.World) {
	i.query.Each(w, func(entry *donburi.Entry) {
		input := component.Input.Get(entry)

		if input.Disabled {
			return
		}

		velocity := component.Velocity.Get(entry)

		velocity.Velocity = math.Vec2{
			X: 0,
			// TODO should match camera scroll speed, get this from settings?
			Y: -0.5,
		}

		touchIDs := ebiten.AppendTouchIDs(nil)
		var isTouch bool
		var touchX, touchY int
		if len(touchIDs) > 0 {
			isTouch = true
			touchX, touchY = ebiten.TouchPosition(touchIDs[0])
		}

		if ebiten.IsKeyPressed(input.MoveUpKey) || (isTouch && touchY < 400) {
			velocity.Velocity.Y = -input.MoveSpeed
		} else if ebiten.IsKeyPressed(input.MoveDownKey) || (isTouch && touchY > 400) {
			velocity.Velocity.Y = input.MoveSpeed
		}

		if ebiten.IsKeyPressed(input.MoveRightKey) || (isTouch && touchX > 200) {
			velocity.Velocity.X = input.MoveSpeed
		}
		if ebiten.IsKeyPressed(input.MoveLeftKey) || (isTouch && touchX < 200) {
			velocity.Velocity.X = -input.MoveSpeed
		}

		// TODO Seems like a very complex way to get the weapon level and timer
		airplane := component.PlayerAirplane.Get(entry)
		player := archetype.MustFindPlayerByNumber(w, airplane.PlayerNumber)
		player.ShootTimer.Update()
		if (ebiten.IsKeyPressed(input.ShootKey) || len(touchIDs) > 1) && player.ShootTimer.IsReady() {
			position := transform.Transform.Get(entry).LocalPosition

			archetype.NewPlayerBullet(w, player, position)
			player.ShootTimer.Reset()
		}
	})
}
