package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
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

type Planet struct {
	radius  float32
	density float32
	mass    float32

	vel     Pair
	pos     Pair
	nextPos Pair
}

type Pair struct {
	X float32
	Y float32
}

func NewPlanet(p Planet) Planet {
	p.mass = math.Pi * p.radius * p.radius * p.density
	p.nextPos = p.pos

	p.vel.X = float32(math.Tan(float64((p.pos.Y-screenHeight/2)/screenHeight))) * 5
	p.vel.Y = float32(math.Tan(float64((p.pos.X-screenWidth/2)/screenWidth))) * -5

	return p
}

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

	// if inpututil.IsKeyJustPressed(ebiten.Key1) {
	// 	g.fps = 10
	// } else if inpututil.IsKeyJustPressed(ebiten.Key2) {
	// 	g.fps = 15
	// } else if inpututil.IsKeyJustPressed(ebiten.Key3) {
	// 	g.fps = 30
	// } else if inpututil.IsKeyJustPressed(ebiten.Key4) {
	// 	g.fps = 60
	// } else if inpututil.IsKeyJustPressed(ebiten.Key5) {
	// 	g.fps = math.MaxInt
	// }

	if g.paused {
		return nil
	}

	// now := time.Now()
	// if now.Sub(g.lastUpdate) < time.Second/g.fps {
	// 	return nil
	// }
	// g.lastUpdate = now

	for i := range g.state {
		g.updatePlanetPosition(i)
	}

	for i := range g.state {
		g.state[i].pos = g.state[i].nextPos
	}

	return nil
}

func (g *Game) updatePlanetPosition(index int) {
	mainPlanet := g.state[index]
	var totalForceX float32 = 0
	var totalForceY float32 = 0
	collision := false

	for i, otherPlanet := range g.state {
		if i == index {
			continue
		}

		xForce, yForce, collision := CalculateForces(mainPlanet, otherPlanet)
		if collision {
			xVel, yVel := collisionVelocities(mainPlanet, otherPlanet)
			g.state[index].vel.X = xVel
			g.state[index].vel.Y = yVel
		}

		totalForceX += xForce
		totalForceY += yForce
	}

	if !collision {
		g.state[index].vel.X += totalForceX / mainPlanet.mass
		g.state[index].vel.Y += totalForceY / mainPlanet.mass

		g.state[index].nextPos.X = g.state[index].pos.X + g.state[index].vel.X
		g.state[index].nextPos.Y = g.state[index].pos.Y + g.state[index].vel.Y
	}
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

func (g *Game) generatePlanets() {
	for i := 0; i < randomPlanetCount; i++ {
		newPlanet := NewPlanet(
			Planet{
				radius:  rand.Float32() * 10,
				density: rand.Float32()*2000 + 15000,

				pos: Pair{
					rand.Float32() * screenHeight,
					rand.Float32() * screenWidth,
				},
			},
		)

		g.state = append(g.state, newPlanet)
	}

	newPlanet := NewPlanet(
		Planet{
			radius:  70,
			density: 2000000000,

			pos: Pair{
				screenHeight / 2,
				screenWidth / 2,
			},
		},
	)
	newPlanet.vel.X = 0
	newPlanet.vel.Y = 0
	g.state = append(g.state, newPlanet)

}

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Conway's Game of Life")

	game := &Game{paused: false, fps: 10, lastUpdate: time.Now()}
	game.generatePlanets()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
