package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
)

type ObserverData struct {
	LookFor *donburi.Query
	Target  *donburi.Entry
}

var Observer = donburi.NewComponentType[ObserverData]()

func ClosestTarget(w donburi.World, entry *donburi.Entry, lookFor *donburi.Query) *donburi.Entry {
	pos := transform.WorldPosition(entry)

	var closestDistance float64
	var closestTarget *donburi.Entry
	lookFor.Each(w, func(target *donburi.Entry) {
		targetPos := transform.WorldPosition(target)
		distance := pos.Distance(targetPos)

		if closestTarget == nil || distance < closestDistance {
			closestTarget = target
			closestDistance = distance
		}
	})

	return closestTarget
}
