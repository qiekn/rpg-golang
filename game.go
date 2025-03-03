package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/qiekn/rpg-go/animations"
	"github.com/qiekn/rpg-go/constants"
	"github.com/qiekn/rpg-go/entities"
	"github.com/qiekn/rpg-go/spritesheet"
)

////////////////////////////////////////////////////////////////////////
//                              Entities                              //
////////////////////////////////////////////////////////////////////////

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	animationFrame    int
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImage      *ebiten.Image
	camera            *Camera
	colliders         []image.Rectangle
}

////////////////////////////////////////////////////////////////////////
//                              Methods                               //
////////////////////////////////////////////////////////////////////////

func NewGame() *Game {

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

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 3,
			Animation: map[entities.PlayerState]*animations.Animation{
				entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
				entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
				entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
				entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
			},
		},
		playerSpriteSheet: playerSpriteSheet,
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
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	return &game
}

func (g *Game) Update() error {

	g.player.Dx = 0.0
	g.player.Dy = 0.0

	// move the player based on keyboard input (wsad)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Dy -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Dy += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Dx -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Dx += 2
	}

	// apply speed & collider check
	g.player.X += g.player.Dx
	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	CheckCollisionVertical(g.player.Sprite, g.colliders)

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

	// enemies AI
	for _, sprite := range g.enemies {
		sprite.Dx = 0.0
		sprite.Dy = 0.0

		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.Dx += 1
			} else if sprite.X > g.player.X {
				sprite.Dx -= 1
			}
			if sprite.Y < g.player.Y {
				sprite.Dy += 1
			} else if sprite.Y > g.player.Y {
				sprite.Dy -= 1
			}
		}

		sprite.X += sprite.Dx
		CheckCollisionHorizontal(sprite.Sprite, g.colliders)

		sprite.Y += sprite.Dy
		CheckCollisionVertical(sprite.Sprite, g.colliders)
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: %d\n", g.player.Health)
		}
	}

	g.camera.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.camera.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
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
			x *= constants.Tilesize
			y *= constants.Tilesize

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))
			// TODO: Bounds? <2025-03-03 13:14, @qiekn> //
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + constants.Tilesize))
			opts.GeoM.Translate(g.camera.X, g.camera.Y)
			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.Frame()
	}

	// draw player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)
	screen.DrawImage(g.player.Img.SubImage(
		g.playerSpriteSheet.Rect(playerFrame),
	).(*ebiten.Image), &opts)
	opts.GeoM.Reset()

	// draw enemies
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}

	// draw potions
	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(potion.Img.SubImage(image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}

	// draw colliders
	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.camera.X),
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{0, 180, 0, 255},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return ebiten.WindowSize()
	return 320, 240
}
