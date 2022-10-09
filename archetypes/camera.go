package archetypes

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewCamera(w donburi.World, startPosition engine.Vector) *donburi.Entry {
	camera := w.Entry(
		w.Create(
			component.Position,
			component.Velocity,
			component.Camera,
		),
	)

	cameraCamera := component.GetCamera(camera)
	cameraCamera.MoveTimer = engine.NewTimer(time.Second * 3)
	donburi.SetValue(camera, component.Position, component.PositionData{
		Position: startPosition,
	})

	return camera
}

func MustFindCamera(w donburi.World) *donburi.Entry {
	camera, ok := query.NewQuery(filter.Contains(component.Camera)).FirstEntity(w)
	if !ok {
		panic("no camera found")
	}

	return camera
}
