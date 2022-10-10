package system

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/archetypes"
	"github.com/m110/airplanes/assets"
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
			filter.Contains(component.Transform, component.Sprite),
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

	cameraPos := component.GetTransform(camera).Position

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
			transform := component.GetTransform(entry)
			sprite := component.GetSprite(entry)

			w, h := sprite.Image.Size()
			halfW, halfH := float64(w)/2, float64(h)/2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-halfW, -halfH)
			op.GeoM.Rotate(float64(int(transform.Rotation)%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(halfW, halfH)

			position := transform.WorldPosition()

			x := position.X
			y := position.Y

			switch sprite.Pivot {
			case component.SpritePivotCenter:
				x -= halfW
				y -= halfH
			}

			op.GeoM.Translate(x, y)

			// TODO Probably not the best place for this?
			// Consider a child object that's hidden sometimes?
			if entry.HasComponent(component.Health) {
				if component.GetHealth(entry).JustDamaged {
					op.ColorM.Translate(1.0, 1.0, 1.0, 0)
				}
			}

			// TODO: Shouldn't be here, consider a child object instead
			if entry.HasComponent(component.PlayerAirplane) {
				airplane := component.GetPlayerAirplane(entry)
				if airplane.Invulnerable {
					op := &ebiten.DrawImageOptions{}
					shieldW, shieldH := assets.AirplaneShield.Size()
					op.GeoM.Translate(x-float64(shieldW)*0.25, y-float64(shieldH)*0.25)
					r.offscreen.DrawImage(assets.AirplaneShield, op)
				}
			}

			r.offscreen.DrawImage(sprite.Image, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cameraPos.X, -cameraPos.Y)
	screen.DrawImage(r.offscreen, op)
}
