package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine/input"
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
		in := component.Input.Get(entry)

		if in.Disabled {
			return
		}

		velocity := component.Velocity.Get(entry)

		velocity.Velocity = math.Vec2{
			X: 0,
			// TODO should match camera scroll speed, get this from settings?
			Y: -0.5,
		}

		if input.IsTouchPrimaryInput() {
			touchIDs := ebiten.AppendTouchIDs(nil)
			var isTouch bool
			var touchX, touchY int
			if len(touchIDs) > 0 {
				isTouch = true
				touchX, touchY = ebiten.TouchPosition(touchIDs[0])
			}

			_ = isTouch
			_ = touchX
			_ = touchY

		} else {
			if ebiten.IsKeyPressed(in.MoveUpKey) {
				velocity.Velocity.Y = -in.MoveSpeed
			} else if ebiten.IsKeyPressed(in.MoveDownKey) {
				velocity.Velocity.Y = in.MoveSpeed
			}

			if ebiten.IsKeyPressed(in.MoveRightKey) {
				velocity.Velocity.X = in.MoveSpeed
			}
			if ebiten.IsKeyPressed(in.MoveLeftKey) {
				velocity.Velocity.X = -in.MoveSpeed
			}
		}

		// TODO Seems like a very complex way to get the weapon level and timer
		airplane := component.PlayerAirplane.Get(entry)
		player := archetype.MustFindPlayerByNumber(w, airplane.PlayerNumber)
		player.ShootTimer.Update()
		if (ebiten.IsKeyPressed(in.ShootKey) || input.IsTouchPrimaryInput()) && player.ShootTimer.IsReady() {
			position := transform.Transform.Get(entry).LocalPosition

			archetype.NewPlayerBullet(w, player, position)
			player.ShootTimer.Reset()
		}
	})
}
