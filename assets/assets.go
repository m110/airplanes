package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"path/filepath"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed ships/ship_0011.png
	shipYellowSmallData []byte
	//go:embed ships/ship_0010.png
	shipGreenSmallData []byte
	//go:embed ships/ship_0018.png
	shipGraySmallData []byte
	//go:embed tiles/tile_0000.png
	laserSingleData []byte

	ShipYellowSmall *ebiten.Image
	ShipGreenSmall  *ebiten.Image
	ShipGraySmall   *ebiten.Image
	LaserSingle     *ebiten.Image

	Levels []Level
)

type Position struct {
	X float64
	Y float64
}

type Level struct {
	Background   *ebiten.Image
	Player1Spawn Position
	Player2Spawn Position
	Enemies      []Enemy
}

type Enemy struct {
	Position Position
	Rotation float64
}

func LoadAssets() {
	ShipYellowSmall = mustNewEbitenImage(shipYellowSmallData)
	ShipGreenSmall = mustNewEbitenImage(shipGreenSmallData)
	ShipGraySmall = mustNewEbitenImage(shipGraySmallData)
	LaserSingle = mustNewEbitenImage(laserSingleData)

	levelPaths, err := filepath.Glob("assets/levels/*.tmx")
	if err != nil {
		panic(err)
	}

	for _, path := range levelPaths {
		Levels = append(Levels, mustLoadLevel(path))
	}
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadLevel(levelPath string) Level {
	levelMap, err := tiled.LoadFile(levelPath)
	if err != nil {
		panic(err)
	}

	level := Level{}

	for _, og := range levelMap.ObjectGroups {
		for _, o := range og.Objects {
			if o.Class == "player1-spawn" {
				level.Player1Spawn = Position{
					X: o.X,
					Y: o.Y,
				}
			}
			if o.Class == "player2-spawn" {
				level.Player2Spawn = Position{
					X: o.X,
					Y: o.Y,
				}
			}
			if o.Class == "enemy-airplane" {
				level.Enemies = append(level.Enemies, Enemy{
					Position: Position{
						X: o.X,
						Y: o.Y,
					},
					Rotation: o.Rotation,
				})
			}
		}
	}

	if level.Player1Spawn == (Position{}) {
		panic("player1-spawn not found")
	}

	if level.Player2Spawn == (Position{}) {
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

	level.Background = ebiten.NewImageFromImage(renderer.Result)

	return level
}
