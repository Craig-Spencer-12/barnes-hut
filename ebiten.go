package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	g.controls()

	if g.paused && !g.stepFrame {
		return nil
	}

	g.BuldTree()
	CalculateCentersOfMassRecursive(g.root)

	g.UpdateAllPlanetPositions()

	g.RemoveDeadPlanets()

	if g.stepFrame {
		g.stepFrame = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.Black)

	g.drawPlanets(screen)

	if g.drawLines {
		g.drawBarnesHutBorders(screen, g.root)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) controls() {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.paused {
		g.stepFrame = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.drawLines = !g.drawLines
		g.BuldTree()
		CalculateCentersOfMassRecursive(g.root)
	}

	g.camera.Controls()
}
