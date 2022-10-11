package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Shooter struct {
	query *query.Query
}

func NewShooter() *Shooter {
	return &Shooter{
		query: query.NewQuery(filter.Contains(component.Transform, component.Shooter)),
	}
}

func (s *Shooter) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		shooter := component.GetShooter(entry)

		shooter.ShootTimer.Update()

		// TODO: It feels like a hack.
		// This relies on another system and requires it to be running before this one.
		// Could be merged into one system, however having them separately also makes sense.
		// Perhaps both components be used by the AI system?
		if entry.HasComponent(component.Observer) {
			observer := component.GetObserver(entry)
			if observer.Target == nil {
				return
			}
		}

		if shooter.ShootTimer.IsReady() {
			shooter.ShootTimer.Reset()

			transform := component.GetTransform(entry)

			archetype.NewEnemyBullet(w, transform.WorldPosition(), transform.WorldRotation())
		}
	})
}
