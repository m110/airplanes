package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/archetype"
	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
	"github.com/m110/airplanes/engine"
	"github.com/m110/airplanes/system"
)

type Airbase struct {
	world     donburi.World
	systems   []System
	drawables []Drawable
}

func NewAirbase(startCallback system.StartGameCallback, backToMenuCallback func()) *Airbase {
	render := system.NewRenderer()

	a := &Airbase{}

	debug := system.NewDebug(a.createWorld)

	a.systems = []System{
		system.NewVelocity(),
		system.NewScript(),
		system.NewPlayerSelect(startCallback, backToMenuCallback),
		debug,
		render,
	}

	a.drawables = []Drawable{
		render,
		system.NewLabel(),
		debug,
	}

	a.createWorld()

	return a
}

func (a *Airbase) createWorld() {
	levelAsset := assets.AirBase
	a.world = donburi.NewWorld()

	levelEntry := a.world.Entry(
		a.world.Create(component.Transform, component.Sprite),
	)

	donburi.SetValue(levelEntry, component.Sprite, component.SpriteData{
		Image: levelAsset.Background,
		Layer: component.SpriteLayerBackground,
		Pivot: component.SpritePivotTopLeft,
	})

	archetype.NewCamera(a.world, engine.Vector{})

	for i, spawn := range levelAsset.Spawns {
		archetype.NewAirbaseAirplane(a.world, spawn.Position, component.MustPlayerFactionFromString(spawn.Faction), i)
	}

	a.world.Create(component.Debug)
}

func (a *Airbase) Update() {
	for _, s := range a.systems {
		s.Update(a.world)
	}
}

func (a *Airbase) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, s := range a.drawables {
		s.Draw(a.world, screen)
	}
}
