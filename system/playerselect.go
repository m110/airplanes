package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type ChosenPlayer struct {
	PlayerNumber int
	Faction      component.PlayerFaction
}

type StartGameCallback func(players []ChosenPlayer)

type PlayerSelect struct {
	query         *query.Query
	startCallback StartGameCallback

	started       bool
	chosenPlayers []ChosenPlayer
}

func NewPlayerSelect(startCallback StartGameCallback) *PlayerSelect {
	return &PlayerSelect{
		query: query.NewQuery(
			filter.Contains(component.Transform, component.PlayerSelect, component.Velocity),
		),
		startCallback: startCallback,
	}
}

func (s *PlayerSelect) Update(w donburi.World) {
	if s.started {
		s.query.EachEntity(w, func(entry *donburi.Entry) {
			playerSelect := component.GetPlayerSelect(entry)
			if !playerSelect.Selected || !playerSelect.Ready {
				return
			}

			transform := component.GetTransform(entry)
			velocity := component.GetVelocity(entry)
			velocity.Velocity.Y -= 0.01

			// TODO: Have something like "find component in children"?
			// Or find by tag
			shadow := transform.Children[1]

			// TODO Better animation
			shadowTransform := component.GetTransform(shadow)
			shadowTransform.LocalPosition.X -= 0.05
			shadowTransform.LocalPosition.Y += 0.05

			// TODO dynamic sprite size
			if transform.WorldPosition().Y <= -32 {
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
				component.GetPlayerSelect(entry).Ready = true
			} else {
				for _, entry := range playerSelects {
					playerSelect := component.GetPlayerSelect(entry)
					if !playerSelect.Selected {
						playerSelect.Selected = true
						playerSelect.PlayerNumber = number
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
								// TODO To methods
								playerSelect.Selected = false
								playerSelect.PlayerNumber = 0

								ps.Selected = true
								ps.PlayerNumber = number
								break
							}
						}
					}
				}
			}
		}

		// TODO Implement
		if inpututil.IsKeyJustPressed(settings.Inputs.Right) {
			if entry, ok := selected[number]; ok {
				playerSelect := component.GetPlayerSelect(entry)
				if !playerSelect.Ready {
					if playerSelect.Index < len(playerSelects)-1 {
						for _, entry := range playerSelects[playerSelect.Index+1:] {
							ps := component.GetPlayerSelect(entry)
							if !ps.Selected {
								// TODO To methods
								playerSelect.Selected = false
								playerSelect.PlayerNumber = 0

								ps.Selected = true
								ps.PlayerNumber = number
								break
							}
						}
					}
				}
			}
		}
	}

	// TODO Fix cancelling ready
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		for i := len(playerSelects) - 1; i >= 0; i-- {
			entry := playerSelects[i]
			playerSelect := component.GetPlayerSelect(entry)
			if playerSelect.Ready {
				playerSelect.Ready = false
				break
			}
			if playerSelect.Selected {
				playerSelect.Selected = false
				playerSelect.PlayerNumber = 0
				break
			}
		}
	}

	for _, entry := range playerSelects {
		playerSelect := component.GetPlayerSelect(entry)
		crosshair := component.GetTransform(entry).Children[0]
		label := component.GetTransform(crosshair).Children[0]

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
