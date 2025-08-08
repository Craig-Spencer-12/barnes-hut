package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Barnes-Hut Planet Sim")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		paused:    startPaused,
		drawLines: startDrawLines,
		camera:    Camera{zoom: 1},
	}

	game.GeneratePlanets()
	game.GenerateSun()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
