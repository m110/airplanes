package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
)

type Progression struct {
	query         *query.Query
	nextLevelFunc func()
}

func NewProgression(nextLevelFunc func()) *Progression {
	return &Progression{
		query: query.NewQuery(
			filter.Contains(
				component.PlayerAirplane,
				component.Velocity,
				component.Input,
				component.Bounds,
			),
		),
		nextLevelFunc: nextLevelFunc,
	}
}

func (p *Progression) Update(w donburi.World) {
	levelEntry := component.MustFindLevel(w)
	level := component.GetLevel(levelEntry)

	if level.Progressed {
		cameraPos := component.GetTransform(archetype.MustFindCamera(w)).LocalPosition
		playersVisible := false
		p.query.EachEntity(w, func(entry *donburi.Entry) {
			playerPos := component.GetTransform(entry).LocalPosition
			playerSprite := component.GetSprite(entry)
			if playerPos.Y+float64(playerSprite.Image.Bounds().Dy()) > cameraPos.Y {
				playersVisible = true
			}
		})
		if !playersVisible {
			p.nextLevelFunc()
		}
		return
	}

	if level.ReachedEnd {
		level.ProgressionTimer.Update()
		if level.ProgressionTimer.IsReady() {
			p.query.EachEntity(w, func(entry *donburi.Entry) {
				input := component.GetInput(entry)
				input.Disabled = true

				velocity := component.GetVelocity(entry)
				velocity.Velocity = engine.Vector{
					X: 0,
					Y: -3,
				}

				bounds := component.GetBounds(entry)
				bounds.Disabled = true
			})

			level.Progressed = true
		}
	} else {
		camera := archetype.MustFindCamera(w)

		cameraPos := component.GetTransform(camera).LocalPosition
		if cameraPos.Y == 0 {
			level.ReachedEnd = true
			level.ProgressionTimer.Reset()
		}
	}
}
