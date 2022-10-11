package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Observer struct {
	query *query.Query
}

func NewObserver() *Observer {
	return &Observer{
		query: query.NewQuery(filter.Contains(component.Transform, component.Observer)),
	}
}

func (s *Observer) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		observer := component.GetObserver(entry)
		if observer.LookFor == nil {
			return
		}

		transform := component.GetTransform(entry)

		var closestDistance float64
		var closestTarget *component.TransformData
		observer.LookFor.EachEntity(w, func(entry *donburi.Entry) {
			targetTransform := component.GetTransform(entry)
			distance := transform.WorldPosition().Distance(targetTransform.WorldPosition())

			if closestTarget == nil || distance < closestDistance {
				closestTarget = targetTransform
				closestDistance = distance
			}
		})

		observer.Target = closestTarget
		if closestTarget == nil {
			return
		}

		// TODO: Should rather rotate towards the target instead of looking at it straight away.
		transform.LookAt(closestTarget.WorldPosition())
	})
}
