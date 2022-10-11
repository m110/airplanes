package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteLayer int

const (
	SpriteLayerBackground SpriteLayer = iota
	SpriteLayerGroundUnits
	SpriteLayerShadows
	SpriteLayerCollectibles
	SpriteLayerAirUnits
	SpriteLayerIndicators
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

	// The original rotation of the sprite
	// "Facing right" is considered 0 degrees
	OriginalRotation float64

	Hidden bool
}

var Sprite = donburi.NewComponentType[SpriteData]()

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
