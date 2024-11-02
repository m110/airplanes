package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Respawn struct {
	query           *donburi.Query
	game            *component.GameData
	restartCallback func()
}

func NewRespawn(restartCallback func()) *Respawn {
	return &Respawn{
		query:           donburi.NewQuery(filter.Contains(component.Player)),
		restartCallback: restartCallback,
	}
}

func (r *Respawn) Update(w donburi.World) {
	if r.game == nil {
		r.game = component.MustFindGame(w)
		if r.game == nil {
			return
		}
	}

	if r.game.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			r.restartCallback()
		}
		return
	}

	playersAlive := 0

	r.query.Each(w, func(entry *donburi.Entry) {
		player := component.Player.Get(entry)

		if player.Lives > 0 {
			playersAlive++
		}

		if player.Respawning {
			player.RespawnTimer.Update()
			if player.RespawnTimer.IsReady() {
				player.Respawning = false
				archetype.NewPlayerAirplane(w, *player, player.PlayerFaction, player.EvolutionLevel())
			}
		}
	})

	// TODO Is this the proper system?
	if playersAlive == 0 {
		game := component.MustFindGame(w)
		if !game.GameOver {
			game.GameOver = true
			cam := archetype.MustFindCamera(w)
			component.Velocity.Get(cam).Velocity.Y = 0
		}
	}
}
