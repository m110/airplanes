package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type DistanceLimit struct {
	query *donburi.Query
}

func NewDistanceLimit() *DistanceLimit {
	return &DistanceLimit{
		query: donburi.NewQuery(
			filter.Contains(component.DistanceLimit),
		),
	}
}

func (l *DistanceLimit) Init(w donburi.World) {}

func (l *DistanceLimit) Update(w donburi.World) {
	l.query.Each(w, func(entry *donburi.Entry) {
		dl := component.DistanceLimit.Get(entry)
		pos := transform.WorldPosition(entry)

		if !dl.Initialized {
			dl.Initialized = true
			dl.PreviousPosition = pos
			return
		}

		distance := pos.Distance(dl.PreviousPosition)
		dl.DistanceTraveled += distance

		if dl.DistanceTraveled > dl.MaxDistance {
			component.Destroy(entry)
		} else {
			dl.PreviousPosition = pos
		}
	})
}
