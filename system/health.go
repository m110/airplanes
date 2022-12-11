package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

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
	h.query.Each(w, func(entry *donburi.Entry) {
		health := component.Health.Get(entry)

		if health.JustDamaged {
			health.DamageIndicatorTimer.Update()
			if health.DamageIndicatorTimer.IsReady() {
				health.HideDamageIndicator()
			}
		} else {
			if health.Health <= 0 {
				EnemyKilledEvent.Publish(w, EnemyKilled{
					Enemy: entry,
				})
			}
		}
	})
}
