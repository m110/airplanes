package system

import (
	"fmt"
	stdmath "math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
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

		var isTouch bool
		var touchX, touchY int
		if input.IsTouchPrimaryInput() {
			touchIDs := ebiten.AppendTouchIDs(nil)
			if len(touchIDs) > 0 {
				isTouch = true
				touchX, touchY = ebiten.TouchPosition(touchIDs[0])
			}
		} else {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				isTouch = true
				touchX, touchY = ebiten.CursorPosition()
			}
		}

		var joystickMovement *math.Vec2
		joystick, joystickFound := engine.FindWithComponent(w, component.Joystick)
		if isTouch {
			if !joystickFound {
				fmt.Println(touchX, touchY)
				joystick = archetype.NewJoystick(w, math.Vec2{
					X: float64(touchX),
					Y: float64(touchY),
				})
			}
			joystickPos := transform.WorldPosition(joystick)
			knobMove := math.Vec2{
				X: float64(touchX) - joystickPos.X,
				Y: float64(touchY) - joystickPos.Y,
			}
			maxRadius := 20.0
			length := knobMove.Magnitude()

			// Get direction vector before clamping (this will be normalized)
			direction := knobMove
			if length > 0 {
				direction = direction.DivScalar(length) // Normalize to length 1
			}

			// Clamp the knob movement
			if length > maxRadius {
				knobMove.X = knobMove.X / length * maxRadius
				knobMove.Y = knobMove.Y / length * maxRadius
			}

			movementFactor := stdmath.Min(length/maxRadius, 1.0)

			knob := hierarchy.MustGetChildren(joystick)[0]
			knobTransform := transform.Transform.Get(knob)
			knobTransform.LocalPosition = knobMove

			// Calculate velocity using the normalized direction vector
			playerSpeed := in.MoveSpeed * movementFactor
			move := direction.MulScalar(playerSpeed)
			joystickMovement = &move
		} else {
			if joystickFound {
				component.Destroy(joystick)
			}
		}

		if in.Disabled {
			return
		}

		velocity := component.Velocity.Get(entry)

		velocity.Velocity = math.Vec2{
			X: 0,
			// TODO should match camera scroll speed, get this from settings?
			Y: -0.5,
		}

		if joystickMovement != nil {
			velocity.Velocity = *joystickMovement
		}

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
