package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

func NewEventBus(w donburi.World) *engine.EventBus {
	count := query.NewQuery(filter.Contains(component.EventBus)).Count(w)
	if count != 0 {
		panic("event bus already exists")
	}

	bus := engine.NewEventBus()
	bus.Subscribe(EnemyKilled{}, OnEnemyKilled)

	eventBus := w.Entry(w.Create(component.EventBus))
	donburi.SetValue(eventBus, component.EventBus, component.EventBusData{
		EventBus: bus,
	})

	return bus
}

type Events struct {
	eventBus *engine.EventBus
}

func NewEvents() *Events {
	return &Events{}
}

func (e *Events) Update(w donburi.World) {
	if e.eventBus == nil {
		e.eventBus = NewEventBus(w)
	}
	e.eventBus.Process(w)
}
