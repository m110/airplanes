package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"

	"github.com/m110/airplanes/assets"
	"github.com/m110/airplanes/component"
)

type HUD struct {
	query         *query.Query
	game          *component.GameData
	shadowOverlay *ebiten.Image
}

func NewHUD() *HUD {
	return &HUD{
		query: query.NewQuery(filter.Contains(component.Player)),
	}
}

func (h *HUD) Draw(w donburi.World, screen *ebiten.Image) {
	if h.game == nil {
		h.game = component.MustFindGame(w)
		if h.game == nil {
			return
		}
		// TODO I don't really like that it's done here
		h.shadowOverlay = ebiten.NewImage(h.game.Settings.ScreenWidth, h.game.Settings.ScreenHeight)
		h.shadowOverlay.Fill(colornames.Black)
	}

	h.query.EachEntity(w, func(entry *donburi.Entry) {
		player := component.GetPlayer(entry)

		icon := assets.Health
		iconWidth, iconHeight := icon.Size()

		baseY := float64(h.game.Settings.ScreenHeight) - float64(iconHeight) - 5
		var baseX float64
		switch player.PlayerNumber {
		case 1:
			baseX = 5
		case 2:
			baseX = float64(h.game.Settings.ScreenWidth) - 5 - float64(iconWidth)
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

	if h.game.GameOver {
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(0, 0, 0, 0.5)
		screen.DrawImage(h.shadowOverlay, op)

		text.Draw(
			screen,
			"GAME OVER",
			assets.NormalFont,
			h.game.Settings.ScreenWidth/4+20,
			h.game.Settings.ScreenHeight/2,
			colornames.White,
		)
	} else if h.game.Paused {
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(0, 0, 0, 0.5)
		screen.DrawImage(h.shadowOverlay, op)

		text.Draw(
			screen,
			"PAUSED",
			assets.NormalFont,
			h.game.Settings.ScreenWidth/4+40,
			h.game.Settings.ScreenHeight/2,
			colornames.White,
		)
	}

	text.Draw(screen, fmt.Sprintf("Score: %06d", h.game.Score), assets.NormalFont, h.game.Settings.ScreenWidth/4, 30, colornames.White)
}
