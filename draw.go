package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawPlanets(screen *ebiten.Image) {
	for _, planet := range g.state {
		x := float32((planet.pos.X - g.camera.pos.X) * g.camera.zoom)
		y := float32((planet.pos.Y - g.camera.pos.Y) * g.camera.zoom)
		r := float32(planet.radius * g.camera.zoom)

		vector.DrawFilledCircle(screen, x, y, r, planet.color, true)
	}
}

func (g *Game) drawBarnesHutBorders(screen *ebiten.Image, n *Node) {
	if n == nil {
		return
	}

	x := float32((n.topLeft.X - g.camera.pos.X) * g.camera.zoom)
	y := float32((n.topLeft.Y - g.camera.pos.Y) * g.camera.zoom)
	size := float32(n.size * g.camera.zoom)

	vector.StrokeLine(screen, x, y, x+size, y, 1, color.White, true)
	vector.StrokeLine(screen, x, y, x, y+size, 1, color.White, true)
	vector.StrokeLine(screen, x+size, y+size, x, y+size, 1, color.White, true)
	vector.StrokeLine(screen, x+size, y+size, x+size, y, 1, color.White, true)

	for i := range n.child {
		g.drawBarnesHutBorders(screen, n.child[i])
	}
}
