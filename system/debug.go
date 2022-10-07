package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/component"
)

type Debug struct {
	debug *component.DebugData
}

func NewDebug() *Debug {
	return &Debug{}
}

func (d *Debug) Update(w donburi.World) {
	if d.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).FirstEntity(w)
		if !ok {
			return
		}

		d.debug = component.GetDebug(debug)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		d.debug.Enabled = !d.debug.Enabled
	}
}

func (d *Debug) Draw(w donburi.World, screen *ebiten.Image) {
	if d.debug == nil || !d.debug.Enabled {
		return
	}

	allCount := w.Len()

	despawnableCount := 0
	spawnedCount := 0
	query.NewQuery(filter.Contains(component.Despawnable)).EachEntity(w, func(entry *donburi.Entry) {
		despawnableCount++
		if component.GetDespawnable(entry).Spawned {
			spawnedCount++
		}
	})

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Entities: %v Despawnable: %v Spawned: %v", allCount, despawnableCount, spawnedCount), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %v", ebiten.CurrentTPS()), 0, 20)
}
