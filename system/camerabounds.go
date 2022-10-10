package system

import (
	"github.com/yohamta/donburi"
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
			component.Transform,
		)),
	}
}

func (b *CameraBounds) Update(w donburi.World) {
	b.query.EachEntity(w, func(entry *donburi.Entry) {
		transform := component.GetTransform(entry)
		if transform.Position.X < 0 {
			transform.Position.X = 0
		}

		if transform.Position.Y < 0 {
			transform.Position.Y = 0
		}
	})
}
