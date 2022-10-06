package component

import "github.com/yohamta/donburi"

const (
	AITypeConstantVelocity = iota
	AITypeFollowPath
)

type AIData struct {
	Spawned bool
	Type    int

	Speed float64

	Path       []PathPosition
	NextTarget int
}

type PathPosition struct {
	X float64
	Y float64
}

var AI = donburi.NewComponentType[AIData]()

func GetAI(entry *donburi.Entry) *AIData {
	return donburi.Get[AIData](entry, AI)
}
