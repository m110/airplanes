package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/lafriks/go-tiled/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

var (
	//go:embed ships/ship_0011.png
	shipYellowSmallData []byte
	//go:embed ships/ship_0010.png
	shipGreenSmallData []byte
	//go:embed tiles/tile_0000.png
	laserSingleData []byte

	ShipYellowSmall *ebiten.Image
	ShipGreenSmall  *ebiten.Image
	LaserSingle     *ebiten.Image

	Level1 *ebiten.Image
)

func LoadAssets() {
	ShipYellowSmall = mustNewEbitenImage(shipYellowSmallData)
	ShipGreenSmall = mustNewEbitenImage(shipGreenSmallData)
	LaserSingle = mustNewEbitenImage(laserSingleData)

	level1, err := tiled.LoadFile("assets/levels/level1.tmx")
	if err != nil {
		panic(err)
	}

	renderer, err := render.NewRenderer(level1)
	if err != nil {
		panic(err)
	}

	err = renderer.RenderLayer(0)
	if err != nil {
		panic(err)
	}

	Level1 = ebiten.NewImageFromImage(renderer.Result)
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
