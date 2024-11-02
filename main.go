package main

import (
	"flag"
	"log"

	"github.com/m110/airplanes/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	quickFlag := flag.Bool("quick", false, "quick mode")
	flag.Parse()

	config := game.Config{
		Quick:        *quickFlag,
		ScreenWidth:  480,
		ScreenHeight: 800,
	}

	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)

	err := ebiten.RunGame(game.NewGame(config))
	if err != nil {
		log.Fatal(err)
	}
}
