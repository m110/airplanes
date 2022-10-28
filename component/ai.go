package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
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

	Path       []math.Vec2
	PathLoops  bool
	NextTarget int
}

var AI = donburi.NewComponentType[AIData]()

func GetAI(entry *donburi.Entry) *AIData {
	return donburi.Get[AIData](entry, AI)
}
