package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type CameraBounds struct {
	query *donburi.Query
}

func NewCameraBounds() *CameraBounds {
	return &CameraBounds{
		query: donburi.NewQuery(filter.Contains(
			component.Camera,
			transform.Transform,
		)),
	}
}

func (b *CameraBounds) Update(w donburi.World) {
	b.query.Each(w, func(entry *donburi.Entry) {
		t := transform.Transform.Get(entry)
		if t.LocalPosition.X < 0 {
			t.LocalPosition.X = 0
		}

		if t.LocalPosition.Y < 0 {
			t.LocalPosition.Y = 0
		}
	})
}
