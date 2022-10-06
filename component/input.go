package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/airplanes/engine"
	"github.com/yohamta/donburi"
)

type InputData struct {
	Disabled bool

	MoveUpKey    ebiten.Key
	MoveRightKey ebiten.Key
	MoveDownKey  ebiten.Key
	MoveLeftKey  ebiten.Key
	MoveSpeed    float64

	ShootKey   ebiten.Key
	ShootTimer *engine.Timer
}

var Input = donburi.NewComponentType[InputData]()

func GetInput(entry *donburi.Entry) *InputData {
	return donburi.Get[InputData](entry, Input)
}
