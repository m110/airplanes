package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/component"
)

type Debug struct {
	query     *query.Query
	debug     *component.DebugData
	offscreen *ebiten.Image
}

func NewDebug() *Debug {
	return &Debug{
		query: query.NewQuery(
			filter.Contains(component.Transform, component.Sprite),
		),
		// TODO figure out the proper size
		offscreen: ebiten.NewImage(3000, 3000),
	}
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

	if d.debug.Enabled {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			query.NewQuery(filter.Contains(component.Player)).EachEntity(w, func(entry *donburi.Entry) {
				player := component.GetPlayer(entry)
				player.UpgradeWeapon()
			})
		}
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

	d.offscreen.Clear()
	d.query.EachEntity(w, func(entry *donburi.Entry) {
		transform := component.GetTransform(entry)
		sprite := component.GetSprite(entry)

		position := transform.WorldPosition()

		w, h := sprite.Image.Size()
		halfW, halfH := float64(w)/2, float64(h)/2

		x := position.X
		y := position.Y

		switch sprite.Pivot {
		case component.SpritePivotCenter:
			x -= halfW
			y -= halfH
		}

		ebitenutil.DrawRect(d.offscreen, transform.Position.X-2, transform.Position.Y-2, 4, 4, colornames.Lime)
		ebitenutil.DebugPrintAt(d.offscreen, fmt.Sprintf("%v", entry.Entity().Id()), int(x), int(y))

		if entry.HasComponent(component.Collider) {
			collider := component.GetCollider(entry)
			ebitenutil.DrawLine(d.offscreen, x, y, x+collider.Width, y, colornames.Lime)
			ebitenutil.DrawLine(d.offscreen, x, y, x, y+collider.Height, colornames.Lime)
			ebitenutil.DrawLine(d.offscreen, x+collider.Width, y, x+collider.Width, y+collider.Height, colornames.Lime)
			ebitenutil.DrawLine(d.offscreen, x, y+collider.Height, x+collider.Width, y+collider.Height, colornames.Lime)
		}

		if entry.HasComponent(component.AI) {
			ai := component.GetAI(entry)
			for i, p := range ai.Path {
				ebitenutil.DrawRect(d.offscreen, p.X-2, p.Y-2, 4, 4, colornames.Red)
				if i < len(ai.Path)-1 {
					next := ai.Path[i+1]
					ebitenutil.DrawLine(d.offscreen, p.X, p.Y, next.X, next.Y, colornames.Red)
				}
			}
		}
	})

	camera := archetypes.MustFindCamera(w)
	cameraPos := component.GetTransform(camera).Position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(d.offscreen, op)
}
