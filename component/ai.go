package component

import "github.com/yohamta/donburi"

const (
	AITypeConstantVelocity = iota
	AITypeFollowPath
)

type AIData struct {
	Spawned bool

	Type             int
	ConstantVelocity float64
	Path             []PositionData
}

var AI = donburi.NewComponentType[AIData]()

func GetAI(entry *donburi.Entry) *AIData {
	return donburi.Get[AIData](entry, AI)
}
