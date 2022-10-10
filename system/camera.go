package system

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Camera struct{}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) Update(w donburi.World) {
	camera := archetypes.MustFindCamera(w)
	cam := component.GetCamera(camera)

	if !cam.Moving {
		cam.MoveTimer.Update()
		if cam.MoveTimer.IsReady() {
			cam.Moving = true
			component.GetVelocity(camera).Velocity.Y = -0.5
		}
	}
}
