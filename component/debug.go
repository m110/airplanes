package component

import "github.com/yohamta/donburi"

type DebugData struct {
	Enabled bool
}

var Debug = donburi.NewComponentType[DebugData]()
