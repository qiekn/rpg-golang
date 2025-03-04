package scenes

import "github.com/hajimehoshi/ebiten/v2"

type SceneId uint

const (
	StartSceneId SceneId = iota
	GameSceneId
	PauseSceneId
	ExitSceneId
)

type Scene interface {
	Start()
	Update() SceneId
	Draw(screen *ebiten.Image)
	OnEnter()
	OnExit()
	IsLoaded() bool
}
