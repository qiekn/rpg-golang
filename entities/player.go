package entities

import "github.com/qiekn/rpg-go/animations"

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Health    uint
	Animation map[PlayerState]*animations.Animation
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animation[Right]
	}
	if dx < 0 {
		return p.Animation[Left]
	}
	if dy > 0 {
		return p.Animation[Down]
	}
	if dy < 0 {
		return p.Animation[Up]
	}
	return nil
}
