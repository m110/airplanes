package archetype

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
			component.Transform,
			component.Velocity,
			component.Camera,
		),
	)

	cameraCamera := component.GetCamera(camera)
	cameraCamera.MoveTimer = engine.NewTimer(time.Second * 3)
	donburi.SetValue(camera, component.Transform, component.TransformData{
		LocalPosition: startPosition,
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
