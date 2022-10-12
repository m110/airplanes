package component

import "github.com/yohamta/donburi"

type LabelData struct {
	Text   string
	Hidden bool
}

var Label = donburi.NewComponentType[LabelData]()

func GetLabel(entry *donburi.Entry) *LabelData {
	return donburi.Get[LabelData](entry, Label)
}
