package spritesheet

import "image"

type SpriteSheet struct {
	Width    int
	Height   int
	Tilesize int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := index % s.Width * s.Tilesize
	y := index / s.Width * s.Tilesize

	return image.Rect(
		x, y, x+s.Tilesize, y+s.Tilesize,
	)
}

func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w, h, t,
	}
}
