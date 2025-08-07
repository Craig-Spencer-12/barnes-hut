package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenHeight = 1024
	screenWidth  = 1024

	randomPlanetCount = 10000
	startPaused       = false
	startDrawLines    = false

	maxPlanetSize = 3
	sunSize       = 70

	collisionCooldown = 5
	collisionsOn      = false

	theta       = 0.5    // tuning parameter for accuracy vs performance
	restitution = 0.6    // how elastic are collisions
	minDistance = 0.0001 // prevents infinite recursion when building the tree
)

type Game struct {
	state []Planet
	root  *Node

	paused     bool
	drawLines  bool
	fps        time.Duration
	lastUpdate time.Time
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.drawLines = !g.drawLines
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.fps = 10
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.fps = 15
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.fps = 30
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.fps = 60
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.fps = math.MaxInt
	}

	g.CreateTree()
	CalculateCentersOfMass(g.root)

	if g.paused {
		return nil
	}

	now := time.Now()
	if now.Sub(g.lastUpdate) < time.Second/g.fps {
		return nil
	}
	g.lastUpdate = now

	g.updateAllPlanets()

	for i := range g.state {
		g.state[i].pos = g.state[i].nextPos
	}

	for i, planet := range g.state {
		if planet.isDead && !planet.isSun {
			g.state = append(g.state[:i], g.state[i+1:]...)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.Black)

	if g.drawLines {
		drawBarnesHutBorders(screen, g.root)
	}

	for _, planet := range g.state {
		vector.DrawFilledCircle(screen, planet.pos.X, planet.pos.Y, planet.radius, planet.color, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Barnes-Hut Planet Sim")

	game := &Game{paused: startPaused, drawLines: startDrawLines, fps: math.MaxInt, lastUpdate: time.Now()}
	game.generatePlanets()
	game.generateSun()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
