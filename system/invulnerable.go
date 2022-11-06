package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Invulnerable struct {
	query *query.Query
}

func NewInvulnerable() *Invulnerable {
	return &Invulnerable{
		query: query.NewQuery(filter.Contains(component.PlayerAirplane)),
	}
}

func (s *Invulnerable) Update(w donburi.World) {
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		player := component.PlayerAirplane.Get(entry)
		if player.Invulnerable {
			player.InvulnerableTimer.Update()
			if player.InvulnerableTimer.IsReady() {
				player.StopInvulnerability()
			}
		}
	})
}
