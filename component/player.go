package component

import (
	"fmt"
	"time"

	"github.com/yohamta/donburi"

	"github.com/m110/airplanes/engine"
)

type PlayerFaction int

const (
	PlayerFactionBlue PlayerFaction = iota
	PlayerFactionRed
	PlayerFactionGreen
	PlayerFactionYellow
)

func MustPlayerFactionFromString(s string) PlayerFaction {
	switch s {
	case "blue":
		return PlayerFactionBlue
	case "red":
		return PlayerFactionRed
	case "green":
		return PlayerFactionGreen
	case "yellow":
		return PlayerFactionYellow
	default:
		panic(fmt.Sprintf("unknown player faction: %v", s))
	}
}

type WeaponLevel int

const (
	WeaponLevelSingle WeaponLevel = iota
	WeaponLevelSingleFast
	WeaponLevelDouble
	WeaponLevelDoubleFast
	WeaponLevelDiagonal
	WeaponLevelDoubleDiagonal
)

type PlayerData struct {
	PlayerNumber  int
	PlayerFaction PlayerFaction

	Lives        int
	Respawning   bool
	RespawnTimer *engine.Timer

	WeaponLevel WeaponLevel
	ShootTimer  *engine.Timer
}

func (d *PlayerData) AddLive() {
	d.Lives++
}

func (d *PlayerData) Damage() {
	if d.Respawning {
		return
	}

	d.Lives--

	if d.Lives > 0 {
		d.Respawning = true
		d.RespawnTimer.Reset()
	}
}

func (d *PlayerData) UpgradeWeapon() {
	if d.WeaponLevel == WeaponLevelDoubleDiagonal {
		return
	}
	d.WeaponLevel++
	d.ShootTimer = engine.NewTimer(d.WeaponCooldown())
}

func (d *PlayerData) WeaponCooldown() time.Duration {
	switch d.WeaponLevel {
	case WeaponLevelSingle:
		return 400 * time.Millisecond
	case WeaponLevelSingleFast:
		return 300 * time.Millisecond
	case WeaponLevelDouble:
		return 300 * time.Millisecond
	case WeaponLevelDoubleFast:
		return 200 * time.Millisecond
	case WeaponLevelDiagonal:
		return 200 * time.Millisecond
	case WeaponLevelDoubleDiagonal:
		return 200 * time.Millisecond
	default:
		panic(fmt.Sprintf("unknown weapon level: %v", d.WeaponLevel))
	}
}

func (d *PlayerData) EvolutionLevel() int {
	switch d.WeaponLevel {
	case WeaponLevelSingle:
		fallthrough
	case WeaponLevelSingleFast:
		return 0
	case WeaponLevelDouble:
		fallthrough
	case WeaponLevelDoubleFast:
		return 1
	case WeaponLevelDiagonal:
		fallthrough
	case WeaponLevelDoubleDiagonal:
		return 2
	default:
		panic(fmt.Sprintf("unknown weapon level: %v", d.WeaponLevel))
	}
}

var Player = donburi.NewComponentType[PlayerData]()
