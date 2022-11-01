package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Velocity struct {
	query *query.Query
}

func NewVelocity() *Velocity {
	return &Velocity{
		query: query.NewQuery(
			filter.Contains(transform.Transform, component.Velocity),
		),
	}
}

func (v *Velocity) Update(w donburi.World) {
	v.query.EachEntity(w, func(entry *donburi.Entry) {
		t := transform.GetTransform(entry)
		velocity := component.GetVelocity(entry)

		t.LocalPosition = t.LocalPosition.Add(velocity.Velocity)
	})
}
