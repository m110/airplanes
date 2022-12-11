package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type CameraBounds struct {
	query *query.Query
}

func NewCameraBounds() *CameraBounds {
	return &CameraBounds{
		query: query.NewQuery(filter.Contains(
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
