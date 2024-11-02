package component

import "github.com/yohamta/donburi"

type JoystickData struct {
}

var Joystick = donburi.NewComponentType[JoystickData]()
