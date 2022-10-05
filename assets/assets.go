package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed ships/ship_0011.png
	shipYellowSmallData []byte
	//go:embed ships/ship_0010.png
	shipGreenSmallData []byte
	//go:embed tiles/tile_0000.png
	laserSingleData []byte

	//go:embed levels/level1.tmx
	level1Data []byte

	ShipYellowSmall *ebiten.Image
	ShipGreenSmall  *ebiten.Image
	LaserSingle     *ebiten.Image

	Level1 Level
)

type Position struct {
	X float64
	Y float64
}

type Level struct {
	Background   *ebiten.Image
	Player1Spawn Position
	Player2Spawn Position
}

func LoadAssets() {
	ShipYellowSmall = mustNewEbitenImage(shipYellowSmallData)
	ShipGreenSmall = mustNewEbitenImage(shipGreenSmallData)
	LaserSingle = mustNewEbitenImage(laserSingleData)

	Level1 = mustLoadLevel(level1Data)
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadLevel(levelData []byte) Level {
	levelMap, err := tiled.LoadReader("assets/levels", bytes.NewReader(levelData))
	if err != nil {
		panic(err)
	}

	var player1Spawn, player2Spawn Position
	for _, og := range levelMap.ObjectGroups {
		for _, o := range og.Objects {
			if o.Class == "player1-spawn" {
				player1Spawn = Position{
					X: o.X,
					Y: o.Y,
				}
			}
			if o.Class == "player2-spawn" {
				player2Spawn = Position{
					X: o.X,
					Y: o.Y,
				}
			}
		}
	}

	if player1Spawn == (Position{}) {
		panic("player1-spawn not found")
	}

	if player2Spawn == (Position{}) {
		panic("player2-spawn not found")
	}

	renderer, err := render.NewRenderer(levelMap)
	if err != nil {
		panic(err)
	}

	err = renderer.RenderLayer(0)
	if err != nil {
		panic(err)
	}

	return Level{
		Background:   ebiten.NewImageFromImage(renderer.Result),
		Player1Spawn: player1Spawn,
		Player2Spawn: player2Spawn,
	}
}
