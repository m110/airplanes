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

	cameraCamera := component.Camera.Get(camera)
	cameraCamera.MoveTimer = engine.NewTimer(time.Second * 3)
	transform.Transform.Get(camera).LocalPosition = startPosition

	return camera
}

func MustFindCamera(w donburi.World) *donburi.Entry {
	camera, ok := query.NewQuery(filter.Contains(component.Camera)).FirstEntity(w)
	if !ok {
		panic("no camera found")
	}

	return camera
}
