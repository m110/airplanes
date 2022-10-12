package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"path/filepath"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/m110/airplanes/engine"
)

var (
	//go:embed tiles/airplane_shield.png
	airplaneShieldData []byte

	//go:embed fonts/kenney-future.ttf
	normalFontData []byte
	//go:embed fonts/kenney-future-narrow.ttf
	narrowFontData []byte

	//go:embed *
	assetsFS embed.FS

	AirplaneYellowSmall *ebiten.Image
	AirplaneGreenSmall  *ebiten.Image
	AirplaneGraySmall   *ebiten.Image

	TankBase *ebiten.Image
	TankGun  *ebiten.Image

	LaserSingle *ebiten.Image
	Rocket      *ebiten.Image

	Health        *ebiten.Image
	WeaponUpgrade *ebiten.Image
	Shield        *ebiten.Image

	AirplaneShield *ebiten.Image

	Levels []Level

	NormalFont font.Face
	NarrowFont font.Face
)

const (
	EnemyClassAirplane = "enemy-airplane"
	EnemyClassTank     = "enemy-tank"

	TilesetClassTiles     = "tiles"
	TilesetClassAirplanes = "airplanes"
)

type Level struct {
	Background *ebiten.Image
	Enemies    []Enemy
}

type Enemy struct {
	Class    string
	Position engine.Vector
	Rotation float64
	Speed    float64
	Path     Path
}

type Path struct {
	Points []engine.Vector
	Loops  bool
}

func MustLoadAssets() {
	loader := newLevelLoader()
	Levels = loader.MustLoadLevels()

	NormalFont = mustLoadFont(normalFontData)
	NarrowFont = mustLoadFont(narrowFontData)

	AirplaneYellowSmall = loader.MustFindTile(TilesetClassAirplanes, "airplane-yellow-small")
	AirplaneGreenSmall = loader.MustFindTile(TilesetClassAirplanes, "airplane-green-small")
	AirplaneGraySmall = loader.MustFindTile(TilesetClassAirplanes, "airplane-gray-small-2")

	TankBase = loader.MustFindTile(TilesetClassTiles, "tank-base")
	TankGun = loader.MustFindTile(TilesetClassTiles, "tank-gun")

	LaserSingle = loader.MustFindTile(TilesetClassTiles, "laser-single")
	Rocket = loader.MustFindTile(TilesetClassTiles, "rocket")

	Health = loader.MustFindTile(TilesetClassTiles, "health")
	WeaponUpgrade = loader.MustFindTile(TilesetClassTiles, "weapon-upgrade")
	Shield = loader.MustFindTile(TilesetClassTiles, "shield")

	AirplaneShield = mustNewEbitenImage(airplaneShieldData)
}

func mustLoadFont(data []byte) font.Face {
	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

type levelLoader struct {
	Tilesets map[string]*tiled.Tileset
}

func newLevelLoader() *levelLoader {
	return &levelLoader{
		Tilesets: make(map[string]*tiled.Tileset),
	}
}

func (l *levelLoader) MustLoadLevels() []Level {
	levelPaths, err := fs.Glob(assetsFS, "levels/*.tmx")
	if err != nil {
		panic(err)
	}

	var levels []Level
	for _, path := range levelPaths {
		levels = append(levels, l.mustLoadLevel(path))
	}

	return levels
}

func (l *levelLoader) mustLoadLevel(levelPath string) Level {
	levelMap, err := tiled.LoadFile(levelPath, tiled.WithFileSystem(assetsFS))
	if err != nil {
		panic(err)
	}

	level := Level{}

	paths := map[uint32]Path{}
	for _, og := range levelMap.ObjectGroups {
		for _, o := range og.Objects {
			if len(o.PolyLines) > 0 {
				var points []engine.Vector
				for _, p := range o.PolyLines {
					for _, pp := range *p.Points {
						points = append(points, engine.Vector{
							X: o.X + pp.X,
							Y: o.Y + pp.Y,
						})
					}
				}
				paths[o.ID] = Path{
					Loops:  false,
					Points: points,
				}
			}
			if len(o.Polygons) > 0 {
				var points []engine.Vector
				for _, p := range o.Polygons {
					for _, pp := range *p.Points {
						points = append(points, engine.Vector{
							X: o.X + pp.X,
							Y: o.Y + pp.Y,
						})
					}
				}
				paths[o.ID] = Path{
					Loops:  true,
					Points: points,
				}
			}
		}
	}

	for _, og := range levelMap.ObjectGroups {
		for _, o := range og.Objects {
			if o.Class == EnemyClassAirplane || o.Class == EnemyClassTank {
				enemy := Enemy{
					Class: o.Class,
					Position: engine.Vector{
						X: o.X,
						Y: o.Y,
					},
					Rotation: o.Rotation,
					Speed:    1,
				}

				for _, p := range o.Properties {
					if p.Name == "path" {
						pathID, err := strconv.Atoi(p.Value)
						if err != nil {
							panic("invalid path id: " + err.Error())
						}

						path, ok := paths[uint32(pathID)]
						if !ok {
							panic("path not found: " + p.Value)
						}

						enemy.Path = path
					}
					if p.Name == "speed" {
						speed, err := strconv.ParseFloat(p.Value, 64)
						if err != nil {
							panic("invalid speed: " + err.Error())
						}

						enemy.Speed = speed
					}
				}

				level.Enemies = append(level.Enemies, enemy)
			}
		}
	}

	renderer, err := render.NewRendererWithFileSystem(levelMap, assetsFS)
	if err != nil {
		panic(err)
	}

	err = renderer.RenderLayer(0)
	if err != nil {
		panic(err)
	}

	level.Background = ebiten.NewImageFromImage(renderer.Result)

	for _, ts := range levelMap.Tilesets {
		if _, ok := l.Tilesets[ts.Class]; !ok {
			l.Tilesets[ts.Class] = ts
		}
	}

	return level
}

func (l *levelLoader) MustFindTile(tilesetClass string, tileClass string) *ebiten.Image {
	ts, ok := l.Tilesets[tilesetClass]
	if !ok {
		panic(fmt.Sprintf("tileset not found: %s", tilesetClass))
	}

	for _, t := range ts.Tiles {
		f, err := fs.ReadFile(assetsFS, filepath.Join("levels", ts.Image.Source))
		if err != nil {
			panic(err)
		}

		tilesetImage := mustNewEbitenImage(f)
		if t.Class == tileClass {
			width := ts.TileWidth
			height := ts.TileHeight

			col := int(t.ID) % ts.Columns
			row := int(t.ID) / ts.Columns

			// Plus one because of 1px margin
			if col > 0 {
				width += 1
			}
			if row > 0 {
				height += 1
			}

			sx := col * width
			sy := row * height

			return tilesetImage.SubImage(
				image.Rect(sx, sy, sx+ts.TileWidth, sy+ts.TileHeight),
			).(*ebiten.Image)
		}
	}

	panic(fmt.Sprintf("tile not found: %s", tileClass))
}
