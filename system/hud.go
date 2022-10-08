package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

type HUD struct {
	query    *query.Query
	settings *component.SettingsData
}

func NewHUD() *HUD {
	return &HUD{
		query: query.NewQuery(filter.Contains(component.Player)),
	}
}

func (h *HUD) Draw(w donburi.World, screen *ebiten.Image) {
	if h.settings == nil {
		h.settings = component.MustFindSettings(w)
		if h.settings == nil {
			return
		}
	}

	h.query.EachEntity(w, func(entry *donburi.Entry) {
		player := component.GetPlayer(entry)

		icon := assets.Health
		iconWidth, iconHeight := icon.Size()

		baseY := float64(h.settings.ScreenHeight) - float64(iconHeight) - 5
		var baseX float64
		switch player.PlayerNumber {
		case 1:
			baseX = 5
		case 2:
			baseX = float64(h.settings.ScreenWidth) - 5 - float64(iconWidth)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(baseX, baseY)
		for i := 0; i < player.Lives; i++ {
			if i > 0 {
				op.GeoM.Translate(0, -float64(iconHeight+2))
			}
			screen.DrawImage(icon, op)
		}
	})
}
