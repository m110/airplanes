package system

import (
	"math/rand"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Health struct {
	query *query.Query
}

func NewHealth() *Health {
	return &Health{
		query: query.NewQuery(filter.Contains(
			transform.Transform,
			component.Health,
		)),
	}
}

func (h *Health) Update(w donburi.World) {
	h.query.EachEntity(w, func(entry *donburi.Entry) {
		health := component.GetHealth(entry)

		if health.JustDamaged {
			health.DamageIndicatorTimer.Update()
			if health.DamageIndicatorTimer.IsReady() {
				health.HideDamageIndicator()
			}
		} else {
			if health.Health <= 0 {
				r := rand.Intn(10)
				if r < 7 {
					archetype.NewRandomCollectible(w, transform.GetTransform(entry).LocalPosition)
				}

				// TODO: It seems like a good candidate to be triggered by an event.
				component.MustFindGame(w).AddScore(1)

				if entry.HasComponent(component.Wreckable) {
					archetype.NewAirplaneWreck(w, entry, component.GetSprite(entry))
				}

				hierarchy.RemoveRecursive(entry)
			}
		}
	})
}
