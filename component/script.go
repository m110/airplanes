package component

import "github.com/yohamta/donburi"

type ScriptData struct {
	Update func(w donburi.World)
}

var Script = donburi.NewComponentType[ScriptData]()

func GetScript(entry *donburi.Entry) *ScriptData {
	return donburi.Get[ScriptData](entry, Script)
}
