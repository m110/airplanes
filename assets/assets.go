package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"io/fs"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/m110/airplanes/engine"
)

var (
	//go:embed airplanes/0011.png
	airplaneYellowSmallData []byte
	//go:embed airplanes/0010.png
	airplaneGreenSmallData []byte
	//go:embed airplanes/0018.png
	airplaneGraySmallData []byte

	//go:embed tiles/tile_0029.png
	tankBaseData []byte
	//go:embed tiles/tile_0030.png
	tankGunData []byte

	//go:embed tiles/tile_0016.png
	turretBaseData []byte
	//go:embed tiles/tile_0018.png
	turretGunSingleData []byte
	//go:embed tiles/tile_0017.png
	turretGunDoubleData []byte

	//go:embed tiles/tile_0000.png
	laserSingleData []byte
	//go:embed tiles/tile_0001.png
	laserDoubleData []byte
	//go:embed tiles/tile_0012.png
	rocketData []byte
	//go:embed tiles/tile_0006.png
	bombData []byte

	//go:embed tiles/tile_0007.png
	flashData []byte
	//go:embed tiles/tile_0008.png
	smokeData []byte

	//go:embed tiles/tile_0024.png
	healthData []byte
	//go:embed tiles/tile_0025.png
	weaponUpgradeData []byte
	//go:embed tiles/tile_0026.png
	shieldData []byte
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

	TurretBase      *ebiten.Image
	TurretGunSingle *ebiten.Image
	TurretGunDouble *ebiten.Image

	LaserSingle *ebiten.Image
	LaserDouble *ebiten.Image
	Rocket      *ebiten.Image
	Bomb        *ebiten.Image

	Smoke *ebiten.Image
	Flash *ebiten.Image

	Health         *ebiten.Image
	WeaponUpgrade  *ebiten.Image
	Shield         *ebiten.Image
	AirplaneShield *ebiten.Image

	Levels []Level

	NormalFont font.Face
	NarrowFont font.Face
)

const (
	EnemyClassAirplane = "enemy-airplane"
	EnemyClassTank     = "enemy-tank"
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
	AirplaneYellowSmall = mustNewEbitenImage(airplaneYellowSmallData)
	AirplaneGreenSmall = mustNewEbitenImage(airplaneGreenSmallData)
	AirplaneGraySmall = mustNewEbitenImage(airplaneGraySmallData)

	TankBase = mustNewEbitenImage(tankBaseData)
	TankGun = mustNewEbitenImage(tankGunData)

	TurretBase = mustNewEbitenImage(turretBaseData)
	TurretGunSingle = mustNewEbitenImage(turretGunSingleData)
	TurretGunDouble = mustNewEbitenImage(turretGunDoubleData)

	LaserSingle = mustNewEbitenImage(laserSingleData)
	LaserDouble = mustNewEbitenImage(laserDoubleData)
	Rocket = mustNewEbitenImage(rocketData)
	Bomb = mustNewEbitenImage(bombData)

	Flash = mustNewEbitenImage(flashData)
	Smoke = mustNewEbitenImage(smokeData)

	Health = mustNewEbitenImage(healthData)
	WeaponUpgrade = mustNewEbitenImage(weaponUpgradeData)
	Shield = mustNewEbitenImage(shieldData)
	AirplaneShield = mustNewEbitenImage(airplaneShieldData)

	levelPaths, err := fs.Glob(assetsFS, "levels/*.tmx")
	if err != nil {
		panic(err)
	}

	for _, path := range levelPaths {
		Levels = append(Levels, mustLoadLevel(path))
	}

	NormalFont = mustLoadFont(normalFontData)
	NarrowFont = mustLoadFont(narrowFontData)
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

func mustLoadLevel(levelPath string) Level {
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

	return level
}
