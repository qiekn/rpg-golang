package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/qiekn/rpg-go/constants"
	"github.com/qiekn/rpg-go/entities"
)

////////////////////////////////////////////////////////////////////////
//                           Main Function                            //
////////////////////////////////////////////////////////////////////////

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

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
