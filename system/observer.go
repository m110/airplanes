package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type Observer struct {
	query *donburi.Query
}

func NewObserver() *Observer {
	return &Observer{
		query: donburi.NewQuery(filter.Contains(transform.Transform, component.Observer)),
	}
}

func (s *Observer) Update(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		observer := component.Observer.Get(entry)
		if observer.LookFor == nil {
			return
		}

		observer.Target = component.ClosestTarget(w, entry, observer.LookFor)
		if observer.Target == nil {
			return
		}

		// TODO: Should rather rotate towards the target instead of looking at it straight away.
		targetPos := transform.WorldPosition(observer.Target)
		transform.LookAt(entry, targetPos)
	})
}
