package component

import "github.com/yohamta/donburi"

type DebugData struct {
	Enabled bool
}

var Debug = donburi.NewComponentType[DebugData]()

func GetDebug(entry *donburi.Entry) *DebugData {
	return donburi.Get[DebugData](entry, Debug)
}
