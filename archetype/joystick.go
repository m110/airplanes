package archetype

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

func NewJoystick(w donburi.World, pos math.Vec2) {
	joystick := w.Entry(w.Create(
		transform.Transform,
		component.Joystick,
		component.Sprite,
	))
	component.Sprite.SetValue(joystick, component.SpriteData{
		Image: assets.JoystickBase,
		Layer: component.SpriteLayerUI,
		Pivot: component.SpritePivotCenter,
	})
	t := transform.Transform.Get(joystick)
	t.LocalPosition = pos
	t.LocalScale = math.Vec2{X: 0.5, Y: 0.5}

	knob := w.Entry(w.Create(
		transform.Transform,
		component.Joystick,
		component.Sprite,
	))

	component.Sprite.SetValue(knob, component.SpriteData{
		Image: assets.JoystickKnob,
		Layer: component.SpriteLayerUI,
		Pivot: component.SpritePivotCenter,
	})
	transform.AppendChild(joystick, knob, false)
}
