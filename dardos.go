package main

import (
	_ "image/png"
	"my_dart_game/scenes"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	mainScene := scenes.NewMainScene()
	pixelgl.Run(mainScene.Run)
}
