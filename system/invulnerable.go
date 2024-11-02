package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/component"
)

type Invulnerable struct {
	query *donburi.Query
}

func NewInvulnerable() *Invulnerable {
	return &Invulnerable{
		query: donburi.NewQuery(filter.Contains(component.PlayerAirplane)),
	}
}

func (s *Invulnerable) Update(w donburi.World) {
	s.query.Each(w, func(entry *donburi.Entry) {
		player := component.PlayerAirplane.Get(entry)
		if player.Invulnerable {
			player.InvulnerableTimer.Update()
			if player.InvulnerableTimer.IsReady() {
				player.StopInvulnerability()
			}
		}
	})
}
