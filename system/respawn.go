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
	query        *query.Query
	screenWidth  int
	screenHeight int
}

func NewRespawn(screenWidth int, screenHeight int) *Respawn {
	return &Respawn{
		query:        query.NewQuery(filter.Contains(component.Player)),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
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

				cameraPos := component.GetPosition(archetypes.MustFindCamera(w))

				// TODO Make this more elegant somehow - move to archetypes?
				switch player.PlayerNumber {
				case 1:
					pos := component.PositionData{
						X: float64(r.screenWidth) * 0.25,
						Y: cameraPos.Y + float64(r.screenHeight)*0.9,
					}
					archetypes.NewPlayerOne(w, pos)
				case 2:
					pos := component.PositionData{
						X: float64(r.screenWidth) * 0.75,
						Y: cameraPos.Y + float64(r.screenHeight)*0.9,
					}
					archetypes.NewPlayerTwo(w, pos)
				}
			}
		}
	})

	if playersAlive == 0 {
		fmt.Println("Game Over")
		// TODO Game Over
	}
}
