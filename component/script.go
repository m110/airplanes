package component

import "github.com/yohamta/donburi"

type ScriptData struct {
	Update func(w donburi.World)
}

var Script = donburi.NewComponentType[ScriptData]()
