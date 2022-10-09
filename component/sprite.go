package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteLayer int

const (
	SpriteLayerBackground SpriteLayer = iota
	SpriteLayerGroundUnits
	SpriteLayerAirUnits
)

type SpritePivot int

const (
	SpritePivotCenter SpritePivot = iota
	SpritePivotTopLeft
)

type SpriteData struct {
	Image *ebiten.Image
	Layer SpriteLayer
	Pivot SpritePivot
}

var Sprite = donburi.NewComponentType[SpriteData]()

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
