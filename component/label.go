package component

import "github.com/yohamta/donburi"

type LabelData struct {
	Text   string
	Hidden bool
}

var Label = donburi.NewComponentType[LabelData]()
