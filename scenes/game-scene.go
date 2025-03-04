package scenes

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/qiekn/rpg-go/animations"
	"github.com/qiekn/rpg-go/camera"
	"github.com/qiekn/rpg-go/constants"
	"github.com/qiekn/rpg-go/entities"
	"github.com/qiekn/rpg-go/spritesheet"
	"github.com/qiekn/rpg-go/tilemap"
	"github.com/qiekn/rpg-go/tileset"
)

type GameScene struct {
	loaded            bool
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *tilemap.TilemapJSON
	tilesets          []tileset.Tileset
	tilemapImage      *ebiten.Image
	camera            *camera.Camera
	colliders         []image.Rectangle
}

func NewGameScene() *GameScene {
	return &GameScene{
		loaded:            false,
		player:            nil,
		playerSpriteSheet: nil,
		enemies:           make([]*entities.Enemy, 0),
		potions:           make([]*entities.Potion, 0),
		tilemapJSON:       nil,
		tilesets:          nil,
		tilemapImage:      nil,
		camera:            nil,
		colliders:         make([]image.Rectangle, 0),
	}
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// fill the screen with a color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	// used to set the translation
	opts := ebiten.DrawImageOptions{}

	// loop over the layers
	for layerIndex, layer := range g.tilemapJSON.Layers {
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

func (g *GameScene) IsLoaded() bool {
	return g.IsLoaded()
}

func (g *GameScene) OnEnter() {
	fmt.Println("enter game scene")
}

func (g *GameScene) OnExit() {
	fmt.Println("exit game scene")
}

func (g *GameScene) Start() {
	// load game assets
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

	tilemapJSON, err := tilemap.NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

	// asign properties
	g.player = &entities.Player{
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
	}
	g.playerSpriteSheet = playerSpriteSheet
	g.enemies = []*entities.Enemy{
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
	}
	g.potions = []*entities.Potion{
		{
			Sprite: &entities.Sprite{
				Img: potionImg,
				X:   210.0,
				Y:   100.0,
			},
			AmtHeal: 1.0,
		},
	}
	g.tilemapJSON = tilemapJSON
	g.tilemapImage = tilemapImg
	g.tilesets = tilesets
	g.camera = camera.NewCamera(0.0, 0.0)
	g.colliders = []image.Rectangle{
		image.Rect(100, 100, 116, 116),
	}
	g.loaded = true
}

func (g *GameScene) Update() SceneId {
	// keymap for scene switch
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		return PauseSceneId
	}

	// actual update method
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

	return GameSceneId
}

var _ Scene = (*GameScene)(nil)

////////////////////////////////////////////////////////////////////////
//                          Helper Functions                          //
////////////////////////////////////////////////////////////////////////

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize),
		) {
			// P.S. the sprite sprite is drawn from the top left corner
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize),
		) {
			// P.S. the sprite sprite is drawn from the top left corner
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
