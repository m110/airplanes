package component

import "github.com/yohamta/donburi"

type ParentData struct {
	Children []*donburi.Entry
}

func (d *ParentData) AppendChild(parent, child *donburi.Entry) {
	d.Children = append(d.Children, child)

	if child.HasComponent(Position) {
		GetPosition(child).SetParent(parent)
	}
}

var Parent = donburi.NewComponentType[ParentData]()

func GetParent(entry *donburi.Entry) *ParentData {
	return donburi.Get[ParentData](entry, Parent)
}
