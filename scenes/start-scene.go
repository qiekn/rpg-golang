package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StartScene struct {
	loaded bool
}

func NewStartScene() *StartScene {
	return &StartScene{
		loaded: false,
	}
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	ebitenutil.DebugPrint(screen, "Press enter to start")
}

func (s *StartScene) IsLoaded() bool {
	return s.IsLoaded()
}

func (s *StartScene) OnEnter() {
	fmt.Println("enter start scene")
}

func (s *StartScene) OnExit() {
	fmt.Println("exit start scene")
}

func (s *StartScene) Start() {
	s.loaded = true
}

func (s *StartScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return StartSceneId
}

var _ Scene = (*StartScene)(nil)
