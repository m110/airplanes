package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/component"
)

type Progression struct {
	query         *donburi.Query
	nextLevelFunc func()
}

func NewProgression(nextLevelFunc func()) *Progression {
	return &Progression{
		query: donburi.NewQuery(
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
	level := component.Level.Get(levelEntry)

	if level.Progressed {
		cameraPos := transform.Transform.Get(archetype.MustFindCamera(w)).LocalPosition
		playersVisible := false
		p.query.Each(w, func(entry *donburi.Entry) {
			playerPos := transform.Transform.Get(entry).LocalPosition
			playerSprite := component.Sprite.Get(entry)
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
			p.query.Each(w, func(entry *donburi.Entry) {
				input := component.Input.Get(entry)
				input.Disabled = true

				velocity := component.Velocity.Get(entry)
				velocity.Velocity = math.Vec2{
					X: 0,
					Y: -3,
				}

				bounds := component.Bounds.Get(entry)
				bounds.Disabled = true
			})

			level.Progressed = true
		}
	} else {
		camera := archetype.MustFindCamera(w)

		cameraPos := transform.Transform.Get(camera).LocalPosition
		if cameraPos.Y == 0 {
			level.ReachedEnd = true
			level.ProgressionTimer.Reset()
		}
	}
}
