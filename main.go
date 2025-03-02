package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	PlayerImage *ebiten.Image
	X, Y        float64
}

func (g *Game) Update() error {

	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.X += 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.X, g.Y)

	// draw out player
	screen.DrawImage(
		g.PlayerImage.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image),
		&opts,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	PlayerImage, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{PlayerImage: PlayerImage, X: 100, Y: 100}); err != nil {
		log.Fatal(err)
	}
}
