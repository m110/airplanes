package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

const maxLevel = 2

type EvolutionData struct {
	Level       int
	Evolving    bool
	GrowTimer   *engine.Timer
	ShrinkTimer *engine.Timer
}

func (e *EvolutionData) Evolve() {
	if e.Level >= maxLevel {
		return
	}

	e.Level++
	e.GrowTimer.Reset()
	e.ShrinkTimer.Reset()
	e.Evolving = true
}

func (e *EvolutionData) StopEvolving() {
	e.Evolving = false
}

var Evolution = donburi.NewComponentType[EvolutionData]()

func GetEvolution(entry *donburi.Entry) *EvolutionData {
	return donburi.Get[EvolutionData](entry, Evolution)
}
