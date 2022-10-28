package archetype

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewCamera(w donburi.World, startPosition math.Vec2) *donburi.Entry {
	camera := w.Entry(
		w.Create(
			transform.Transform,
			component.Velocity,
			component.Camera,
		),
	)

	cameraCamera := component.GetCamera(camera)
	cameraCamera.MoveTimer = engine.NewTimer(time.Second * 3)
	donburi.SetValue(camera, transform.Transform, transform.TransformData{
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
