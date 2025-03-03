package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/qiekn/rpg-go/entities"
)

////////////////////////////////////////////////////////////////////////
//                              Entities                              //
////////////////////////////////////////////////////////////////////////

type Game struct {
	player       *entities.Player
	enemies      []*entities.Enemy
	potions      []*entities.Potion
	tilemapJSON  *TilemapJSON
	tilesets     []Tileset
	tilemapImage *ebiten.Image
	camera       *Camera
}

////////////////////////////////////////////////////////////////////////
//                              Methods                               //
////////////////////////////////////////////////////////////////////////

func (g *Game) Update() error {

	// move the player based on keyboard input (wsad)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.X += 1
			} else if sprite.X > g.player.X {
				sprite.X -= 1
			}
			if sprite.Y < g.player.Y {
				sprite.Y += 1
			} else if sprite.Y > g.player.Y {
				sprite.Y -= 1
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: %d\n", g.player.Health)
		}
	}

	g.camera.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.camera.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Height)*16.0,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// fill the screen with a color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	// used to set the translation
	opts := ebiten.DrawImageOptions{}

	for layerIndex, layer := range g.tilemapJSON.Layers {
		// loop over the layers
		// loop over the tiles in the layer
		for index, id := range layer.Data {

			// skip empty tiles (ID 0 represents an empty tile in the tilemap)
			if id == 0 {
				continue
			}

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile postion to pixel position
			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))
			// TODO: Bounds? <2025-03-03 13:14, @qiekn> //
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))
			opts.GeoM.Translate(g.camera.X, g.camera.Y)
			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	// draw player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)
	screen.DrawImage(g.player.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
	opts.GeoM.Reset()

	// draw enemies
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}

	// draw potions
	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(potion.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return ebiten.WindowSize()
	return 320, 240
}

////////////////////////////////////////////////////////////////////////
//                           Main Function                            //
////////////////////////////////////////////////////////////////////////

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// load the image from file
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/tileset-floor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 3,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   50.0,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				AmtHeal: 1.0,
			},
		},
		tilemapJSON:  tilemapJSON,
		tilemapImage: tilemapImg,
		tilesets:     tilesets,
		camera:       NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
