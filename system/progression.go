package system

import (
	"time"

	"github.com/m110/airplanes/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Progression struct {
	query            *query.Query
	progressionTimer *engine.Timer
	reachedEnd       bool
	progressed       bool
}

func NewProgression() *Progression {
	return &Progression{
		query: query.NewQuery(
			filter.Contains(
				component.PlayerTag,
				component.Velocity,
				component.Input,
				component.Bounds,
			),
		),
		progressionTimer: engine.NewTimer(time.Second * 3),
	}
}

func (p *Progression) Update(w donburi.World) {
	if p.progressed {
		// TODO next level
		return
	}

	if p.reachedEnd {
		p.progressionTimer.Update()
		if p.progressionTimer.IsReady() {
			p.query.EachEntity(w, func(entry *donburi.Entry) {
				input := component.GetInput(entry)
				input.Disabled = true

				velocity := component.GetVelocity(entry)
				velocity.X = 0
				velocity.Y = -3

				bounds := component.GetBounds(entry)
				bounds.Disabled = true
			})
			p.progressed = true
		}
	} else {
		camera, ok := query.NewQuery(filter.Contains(component.CameraTag)).FirstEntity(w)
		if !ok {
			panic("no camera found")
		}

		cameraPos := component.GetPosition(camera)
		if cameraPos.Y == 0 {
			p.reachedEnd = true
			p.progressionTimer.Reset()
		}
	}
}
