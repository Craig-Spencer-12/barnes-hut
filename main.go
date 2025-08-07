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
	screenWidth  = 2048

	randomPlanetCount = 1000
	startPaused       = false
	startDrawLines    = false

	maxPlanetSize = 10
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

	camera     Camera
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

	g.camera.Update()

	if g.paused {
		return nil
	}

	g.CreateTree()
	CalculateCentersOfMass(g.root)

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
		g.drawBarnesHutBorders(screen, g.root)
	}

	for _, planet := range g.state {
		x := float32((planet.pos.X - g.camera.pos.X) * g.camera.zoom)
		y := float32((planet.pos.Y - g.camera.pos.Y) * g.camera.zoom)
		r := float32(planet.radius * g.camera.zoom)

		vector.DrawFilledCircle(screen, x, y, r, planet.color, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Barnes-Hut Planet Sim")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		paused:     startPaused,
		drawLines:  startDrawLines,
		fps:        math.MaxInt,
		lastUpdate: time.Now(),
		camera:     Camera{zoom: 1},
	}
	game.generatePlanets()
	game.generateSun()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
