package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type InputData struct {
	MoveUpKey    ebiten.Key
	MoveRightKey ebiten.Key
	MoveDownKey  ebiten.Key
	MoveLeftKey  ebiten.Key
	MoveSpeed    float64

	ShootKey ebiten.Key
}

var Input = donburi.NewComponentType[InputData]()

func GetInput(entry *donburi.Entry) *InputData {
	return donburi.Get[InputData](entry, Input)
}
