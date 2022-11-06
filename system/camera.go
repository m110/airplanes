package system

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Camera struct{}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) Update(w donburi.World) {
	camera := archetype.MustFindCamera(w)
	cam := component.Camera.Get(camera)

	if !cam.Moving {
		cam.MoveTimer.Update()
		if cam.MoveTimer.IsReady() {
			cam.Moving = true
			component.Velocity.Get(camera).Velocity.Y = -0.5
		}
	}
}
