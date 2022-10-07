package system

import (
	"fmt"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Render struct {
	query     *query.Query
	offscreen *ebiten.Image
	debug     *component.DebugData
}

func NewRenderer() *Render {
	return &Render{
		query: query.NewQuery(
			filter.Contains(component.Position, component.Sprite),
		),
		// TODO figure out the proper size
		offscreen: ebiten.NewImage(3000, 3000),
	}
}

func (r *Render) Update(w donburi.World) {
	if r.debug == nil {
		debug, ok := query.NewQuery(filter.Contains(component.Debug)).FirstEntity(w)
		if !ok {
			return
		}

		r.debug = component.GetDebug(debug)
	}
}

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	camera := archetypes.MustFindCamera(w)
	cam := component.GetCamera(camera)

	if !cam.Moving {
		cam.MoveTimer.Update()
		if cam.MoveTimer.IsReady() {
			cam.Moving = true
			component.GetVelocity(camera).Y = -0.5
		}
	}

	cameraPos := component.GetPosition(camera)

	r.offscreen.Clear()

	var entries []*donburi.Entry
	r.query.EachEntity(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := lo.GroupBy(entries, func(entry *donburi.Entry) int {
		return int(component.GetSprite(entry).Layer)
	})
	layers := lo.Keys(byLayer)
	sort.Ints(layers)

	for _, layer := range layers {
		for _, entry := range byLayer[layer] {
			position := component.GetPosition(entry)
			sprite := component.GetSprite(entry)

			op := &ebiten.DrawImageOptions{}
			if entry.HasComponent(component.Rotation) {
				w, h := sprite.Image.Size()
				op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
				angle := component.GetRotation(entry).Angle
				op.GeoM.Rotate(float64(int(angle)%360) * 2 * math.Pi / 360)
			}
			op.GeoM.Translate(position.X, position.Y)
			r.offscreen.DrawImage(sprite.Image, op)

			if r.debug != nil && r.debug.Enabled {
				ebitenutil.DebugPrintAt(r.offscreen, fmt.Sprintf("%v", entry.Entity().Id()), int(position.X-5), int(position.Y-5))
				if entry.HasComponent(component.Collider) {
					collider := component.GetCollider(entry)
					ebitenutil.DrawLine(r.offscreen, position.X, position.Y, position.X+collider.Width, position.Y, colornames.Lime)
					ebitenutil.DrawLine(r.offscreen, position.X, position.Y, position.X, position.Y+collider.Height, colornames.Lime)
					ebitenutil.DrawLine(r.offscreen, position.X+collider.Width, position.Y, position.X+collider.Width, position.Y+collider.Height, colornames.Lime)
					ebitenutil.DrawLine(r.offscreen, position.X, position.Y+collider.Height, position.X+collider.Width, position.Y+collider.Height, colornames.Lime)
				}
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}
