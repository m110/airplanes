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
	SpriteLayerFallingWrecks
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

	ColorOverride *ColorOverride
}

type ColorOverride struct {
	R, G, B, A float64
}

func (s *SpriteData) Show() {
	s.Hidden = false
}

func (s *SpriteData) Hide() {
	s.Hidden = true
}

var Sprite = donburi.NewComponentType[SpriteData]()

func GetSprite(entry *donburi.Entry) *SpriteData {
	return donburi.Get[SpriteData](entry, Sprite)
}
