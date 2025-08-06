package main

import (
	"math"
	"math/rand"
)

func NewPlanet(p Planet) Planet {
	p.mass = math.Pi * p.radius * p.radius * p.density
	p.nextPos = p.pos

	p.vel.X = float32(math.Tan(float64((p.pos.Y-screenHeight/2)/screenHeight))) * 5
	p.vel.Y = float32(math.Tan(float64((p.pos.X-screenWidth/2)/screenWidth))) * -5

	return p
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
}

func (g *Game) generateSun() {
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
