package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const (
	SpriteLayerBackground SpriteLayer = iota
	SpriteLayerUnits
)

type SpriteLayer int

type SpriteData struct {
	Image *ebiten.Image
	Layer SpriteLayer
}

var Sprite = donburi.NewComponentType[SpriteData]()

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
