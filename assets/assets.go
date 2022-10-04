package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed ships/ship_0011.png
	shipYellowSmallData []byte

	ShipYellowSmall *ebiten.Image
)

func LoadAssets() error {
	shipYellowSmallImage, _, err := image.Decode(bytes.NewReader(shipYellowSmallData))
	if err != nil {
		return err
	}

	ShipYellowSmall = ebiten.NewImageFromImage(shipYellowSmallImage)

	return nil
}
