package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

const (
	AITypeConstantVelocity AIType = iota
	AITypeFollowPath
)

type AIType int

type AIData struct {
	Spawned bool
	Type    AIType

	Speed float64

	Path       []engine.Vector
	PathLoops  bool
	NextTarget int
}

var AI = donburi.NewComponentType[AIData]()

func GetAI(entry *donburi.Entry) *AIData {
	return donburi.Get[AIData](entry, AI)
}
