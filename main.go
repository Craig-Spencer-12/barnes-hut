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

	randomPlanetCount = 200
)

type Game struct {
	state []Planet
	tree  any

	paused     bool
	fps        time.Duration
	lastUpdate time.Time
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
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

	if g.paused {
		return nil
	}

	now := time.Now()
	if now.Sub(g.lastUpdate) < time.Second/g.fps {
		return nil
	}
	g.lastUpdate = now

	for i := range g.state {
		g.updatePlanetPosition(i)
	}

	for i := range g.state {
		g.state[i].pos = g.state[i].nextPos
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.Black)

	for _, planet := range g.state {
		vector.DrawFilledCircle(screen, planet.pos.X, planet.pos.Y, planet.radius, color.White, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Conway's Game of Life")

	game := &Game{paused: false, fps: math.MaxInt, lastUpdate: time.Now()}
	game.generatePlanets()
	game.generateSun()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
