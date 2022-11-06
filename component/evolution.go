package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

const maxLevel = 2

type EvolutionData struct {
	Level           int
	Evolving        bool
	StartedEvolving bool
	GrowTimer       *engine.Timer
	ShrinkTimer     *engine.Timer
}

func (e *EvolutionData) Evolve() {
	if e.Level >= maxLevel {
		return
	}

	e.Level++
	e.GrowTimer.Reset()
	e.ShrinkTimer.Reset()
	e.Evolving = true
	e.StartedEvolving = false
}

func (e *EvolutionData) StopEvolving() {
	e.Evolving = false
}

var Evolution = donburi.NewComponentType[EvolutionData]()
