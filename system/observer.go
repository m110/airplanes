package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Observer struct {
	query *query.Query
}

func NewObserver() *Observer {
	return &Observer{
		query: query.NewQuery(filter.Contains(transform.Transform, component.Observer)),
	}
}

func (s *Observer) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		observer := component.GetObserver(entry)
		if observer.LookFor == nil {
			return
		}

		pos := transform.WorldPosition(entry)

		var closestDistance float64
		var closestTarget *donburi.Entry
		observer.LookFor.EachEntity(w, func(target *donburi.Entry) {
			targetPos := transform.WorldPosition(target)
			distance := pos.Distance(&targetPos)

			if closestTarget == nil || distance < closestDistance {
				closestTarget = target
				closestDistance = distance
			}
		})

		observer.Target = closestTarget
		if closestTarget == nil {
			return
		}

		// TODO: Should rather rotate towards the target instead of looking at it straight away.
		targetPos := transform.WorldPosition(closestTarget)
		transform.LookAt(entry, targetPos)
	})
}
