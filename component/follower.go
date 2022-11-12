package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type FollowerData struct {
	Target         *donburi.Entry
	FollowingSpeed float64
	FollowingTimer *engine.Timer
}

var Follower = donburi.NewComponentType[FollowerData]()
