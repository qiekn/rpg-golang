package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	loaded bool
}

func NewPauseScene() *PauseScene {
	return &PauseScene{
		loaded: false,
	}
}

func (p *PauseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	ebitenutil.DebugPrint(screen, "Press p again to unpause")
}

func (p *PauseScene) IsLoaded() bool {
	return p.IsLoaded()
}

func (p *PauseScene) OnEnter() {
	fmt.Println("enter pause scene")
}

func (p *PauseScene) OnExit() {
	fmt.Println("exit pause scene")
}

func (p *PauseScene) Start() {
	p.loaded = true
}

func (p *PauseScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		return GameSceneId
	}
	return PauseSceneId
}

var _ Scene = (*PauseScene)(nil)
