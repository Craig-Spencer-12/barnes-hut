package main

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
)

func NewPlanet(p Planet) Planet {
	p.mass = math.Pi * p.radius * p.radius * p.density
	p.nextPos = p.pos

	p.vel.X = float64(math.Tan(float64((p.pos.Y-screenHeight/2)/screenHeight))) * 5
	p.vel.Y = float64(math.Tan(float64((p.pos.X-screenWidth/2)/screenWidth))) * -5

	if p.color.R == 0 {
		p.color = planetColors[rand.Intn(len(planetColors))]
	}

	return p
}

func (g *Game) updatePlanetPositionSlow(index int) {
	mainPlanet := g.state[index]
	var totalForceX float64 = 0
	var totalForceY float64 = 0
	collision := false

	for i, otherPlanet := range g.state {
		if i == index {
			continue
		}

		xForce, yForce, collision := CalculateForcesSlow(mainPlanet, otherPlanet)
		if collision {
			if otherPlanet.isSun {
				g.state[index].isDead = true
				return
			}

			if collisionsOn {
				xVel, yVel := collisionVelocities(mainPlanet, otherPlanet)
				g.state[index].vel.X = xVel
				g.state[index].vel.Y = yVel
			}
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

func (g *Game) generatePlanets() {
	for i := 0; i < randomPlanetCount; i++ {
		newPlanet := NewPlanet(
			Planet{
				radius:  rand.Float64() * maxPlanetSize,
				density: rand.Float64()*2000 + 15000,

				pos: Pair{
					rand.Float64() * screenWidth,
					rand.Float64() * screenHeight,
				},
			},
		)

		g.state = append(g.state, newPlanet)
	}
}

func (g *Game) generateSun() {
	newPlanet := NewPlanet(
		Planet{
			radius:  sunSize,
			density: 2000000000,
			isSun:   true,

			color: color.RGBA{R: 255, G: 215, B: 0, A: 255},

			pos: Pair{
				screenWidth / 2,
				screenHeight / 2,
			},
		},
	)
	newPlanet.vel.X = 0
	newPlanet.vel.Y = 0
	g.state = append(g.state, newPlanet)
}

func (g *Game) updateAllPlanets() {

	var wg sync.WaitGroup
	for i := range g.state {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			g.updatePlanetPosition(i)
			if g.state[i].collisionCooldown > 0 {
				g.state[i].collisionCooldown--
			}
		}(i)
	}
	wg.Wait()
}

func movingTogether(planetA, planetB Planet) bool {
	// Relative velocity
	dvx := planetB.vel.X - planetA.vel.X
	dvy := planetB.vel.Y - planetA.vel.Y

	// Relative position
	dx := planetB.pos.X - planetA.pos.X
	dy := planetB.pos.Y - planetA.pos.Y

	// Dot product to check if moving toward each other
	dot := dvx*dx + dvy*dy

	if dot < 0 {
		return false
	}

	return true
}
