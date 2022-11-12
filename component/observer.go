package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/query"
)

type ObserverData struct {
	LookFor *query.Query
	Target  *donburi.Entry
}

var Observer = donburi.NewComponentType[ObserverData]()

func ClosestTarget(w donburi.World, entry *donburi.Entry, lookFor *query.Query) *donburi.Entry {
	pos := transform.WorldPosition(entry)

	var closestDistance float64
	var closestTarget *donburi.Entry
	lookFor.EachEntity(w, func(target *donburi.Entry) {
		targetPos := transform.WorldPosition(target)
		distance := pos.Distance(targetPos)

		if closestTarget == nil || distance < closestDistance {
			closestTarget = target
			closestDistance = distance
		}
	})

	return closestTarget
}
