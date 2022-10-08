package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type SettingsData struct {
	ScreenWidth  int
	ScreenHeight int
}

var Settings = donburi.NewComponentType[SettingsData]()

func GetSettings(entry *donburi.Entry) *SettingsData {
	return donburi.Get[SettingsData](entry, Settings)
}

func MustFindSettings(w donburi.World) *SettingsData {
	settings, ok := query.NewQuery(filter.Contains(Settings)).FirstEntity(w)
	if !ok {
		panic("settings not found")
	}

	return GetSettings(settings)
}
