package system

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Respawn struct {
	query *query.Query
}

func NewRespawn() *Respawn {
	return &Respawn{
		query: query.NewQuery(filter.Contains(component.Player)),
	}
}

func (r *Respawn) Update(w donburi.World) {
	playersAlive := 0

	r.query.EachEntity(w, func(entry *donburi.Entry) {
		player := component.GetPlayer(entry)

		if player.Lives > 0 {
			playersAlive++
		}

		if player.Respawning {
			player.RespawnTimer.Update()
			if player.RespawnTimer.IsReady() {
				player.Respawning = false
				archetypes.NewPlayerAirplane(w, player.PlayerNumber)
			}
		}
	})

	if playersAlive == 0 {
		fmt.Println("Game Over")
		// TODO Game Over
	}
}
