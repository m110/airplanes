package component

import (
	"github.com/yohamta/donburi"
)

const (
	CollisionLayerPlayerBullets ColliderLayer = iota
	CollisionLayerEnemyBullets
	CollisionLayerGroundEnemies
	CollisionLayerAirEnemies
	CollisionLayerPlayers
	CollisionLayerCollectibles
)

type ColliderLayer int

type ColliderData struct {
	Width  float64
	Height float64
	Layer  ColliderLayer
}

var Collider = donburi.NewComponentType[ColliderData]()
