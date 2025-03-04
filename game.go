package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/qiekn/rpg-go/scenes"
)

type Game struct {
	sceneMap      map[SceneId]Scene
	activeSceneId SceneId
}

func NewGame() *Game {
	sceneMap := map[SceneId]Scene{
		StartSceneId: NewStartScene(),
		GameSceneId:  NewGameScene(),
		PauseSceneId: NewPauseScene(),
	}
	activeSceneId := StartSceneId

	sceneMap[activeSceneId].Start()
	return &Game{
		sceneMap,
		activeSceneId,
	}
}

func (g *Game) Update() error {
	nextSceneId := g.sceneMap[g.activeSceneId].Update()

	// quit application
	if nextSceneId == ExitSceneId {
		g.sceneMap[g.activeSceneId].OnExit()
		return ebiten.Termination
	}

	// switched scenes
	if nextSceneId != g.activeSceneId {
		nextScene := g.sceneMap[nextSceneId]
		g.sceneMap[g.activeSceneId].OnExit()
		nextScene.Start()
		nextScene.OnEnter()
	}

	g.activeSceneId = nextSceneId
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
