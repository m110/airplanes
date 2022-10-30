package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/engine"
)

type EventBusData struct {
	EventBus *engine.EventBus
}

var EventBus = donburi.NewComponentType[EventBusData]()

func GetEventBus(entry *donburi.Entry) *EventBusData {
	return donburi.Get[EventBusData](entry, EventBus)
}

func MustFindEventBus(w donburi.World) *engine.EventBus {
	eventBus, ok := query.NewQuery(filter.Contains(EventBus)).FirstEntity(w)
	if !ok {
		panic("event bus not found")
	}
	return GetEventBus(eventBus).EventBus
}
