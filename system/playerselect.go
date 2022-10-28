package system

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type ChosenPlayer struct {
	PlayerNumber int
	Faction      component.PlayerFaction
}

type StartGameCallback func(players []ChosenPlayer)

type PlayerSelect struct {
	query              *query.Query
	startCallback      StartGameCallback
	backToMenuCallback func()

	started       bool
	altitudeTimer *engine.Timer
	chosenPlayers []ChosenPlayer
}

func NewPlayerSelect(startCallback StartGameCallback, backToMenuCallback func()) *PlayerSelect {
	return &PlayerSelect{
		query: query.NewQuery(
			filter.Contains(
				transform.Transform,
				component.PlayerSelect,
				component.Velocity,
				component.Altitude,
			),
		),
		startCallback:      startCallback,
		backToMenuCallback: backToMenuCallback,
		altitudeTimer:      engine.NewTimer(time.Second),
	}
}

func (s *PlayerSelect) Update(w donburi.World) {
	if s.started {
		s.query.EachEntity(w, func(entry *donburi.Entry) {
			playerSelect := component.GetPlayerSelect(entry)
			if !playerSelect.Selected || !playerSelect.Ready {
				return
			}

			velocity := component.GetVelocity(entry)
			velocity.Velocity.Y -= 0.01

			s.altitudeTimer.Update()
			if s.altitudeTimer.IsReady() {
				component.GetAltitude(entry).Velocity = 0.005
			}

			// TODO Scale the sprite slightly

			// TODO dynamic sprite size not hardcoded
			if transform.WorldPosition(entry).Y <= -32 {
				s.startCallback(s.chosenPlayers)
			}
		})

		return
	}

	var playerSelects []*donburi.Entry
	selected := map[int]*donburi.Entry{}
	s.query.EachEntity(w, func(entry *donburi.Entry) {
		playerSelect := component.GetPlayerSelect(entry)
		if playerSelect.Selected {
			selected[playerSelect.PlayerNumber] = entry
		}

		playerSelects = append(playerSelects, entry)
	})

	for number, settings := range archetype.Players {
		if inpututil.IsKeyJustPressed(settings.Inputs.Shoot) {
			if entry, ok := selected[number]; ok {
				component.GetPlayerSelect(entry).LockIn()
			} else {
				for _, entry := range playerSelects {
					playerSelect := component.GetPlayerSelect(entry)
					if !playerSelect.Selected {
						playerSelect.Select(number)
						break
					}
				}
			}
		}

		if inpututil.IsKeyJustPressed(settings.Inputs.Left) {
			if entry, ok := selected[number]; ok {
				playerSelect := component.GetPlayerSelect(entry)
				if !playerSelect.Ready {
					// TODO refactor
					if playerSelect.Index > 0 {
						for i := playerSelect.Index - 1; i >= 0; i-- {
							entry := playerSelects[i]
							ps := component.GetPlayerSelect(entry)
							if !ps.Selected {
								playerSelect.Unselect()

								ps.Select(number)
								break
							}
						}
					}
				}
			}
		}

		if inpututil.IsKeyJustPressed(settings.Inputs.Right) {
			if entry, ok := selected[number]; ok {
				playerSelect := component.GetPlayerSelect(entry)
				if !playerSelect.Ready {
					if playerSelect.Index < len(playerSelects)-1 {
						for _, entry := range playerSelects[playerSelect.Index+1:] {
							ps := component.GetPlayerSelect(entry)
							if !ps.Selected {
								playerSelect.Unselect()

								ps.Select(number)
								break
							}
						}
					}
				}
			}
		}
	}

	// TODO Cancel just the last action
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cancelled := false
		for _, entry := range playerSelects {
			playerSelect := component.GetPlayerSelect(entry)
			if playerSelect.Ready {
				playerSelect.Release()
				cancelled = true
			}
			if playerSelect.Selected {
				playerSelect.Unselect()
				cancelled = true
			}
		}

		if !cancelled {
			s.backToMenuCallback()
		}
	}

	for _, entry := range playerSelects {
		playerSelect := component.GetPlayerSelect(entry)
		crosshair := w.Entry(hierarchy.MustGetChildren(entry)[0])
		label := w.Entry(hierarchy.MustGetChildren(crosshair)[0])

		if playerSelect.Selected {
			if playerSelect.Ready {
				component.GetSprite(crosshair).Hidden = true
			} else {
				component.GetSprite(crosshair).Hidden = false
			}

			component.GetLabel(label).Text = fmt.Sprintf("Player %v", playerSelect.PlayerNumber)
			component.GetLabel(label).Hidden = false
		} else {
			component.GetSprite(crosshair).Hidden = true
			component.GetLabel(label).Hidden = true
		}
	}

	var chosenPlayers []ChosenPlayer
	playersReady := 0
	for _, entry := range playerSelects {
		playerSelect := component.GetPlayerSelect(entry)
		if playerSelect.Selected {
			chosenPlayers = append(chosenPlayers, ChosenPlayer{
				PlayerNumber: playerSelect.PlayerNumber,
				Faction:      playerSelect.Faction,
			})
			if playerSelect.Ready {
				playersReady++
			}
		}
	}

	if playersReady > 0 && playersReady == len(chosenPlayers) {
		s.chosenPlayers = chosenPlayers
		s.started = true
		for _, entry := range playerSelects {
			ps := component.GetPlayerSelect(entry)
			if ps.Selected && ps.Ready {
				velocity := component.GetVelocity(entry)
				velocity.Velocity.Y = -0.5
			}
		}
	}
}
