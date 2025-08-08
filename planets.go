package main

import (
	"math"
	"math/rand"
	"sync"
)

func NewPlanet(p Planet) Planet {
	p.mass = math.Pi * p.radius * p.radius * p.density
	p.nextPos = p.pos

	p.vel.X = math.Tan((p.pos.Y-screenHeight/2)/screenHeight) * startingVelocity
	p.vel.Y = math.Tan((p.pos.X-screenWidth/2)/screenWidth) * -startingVelocity

	if p.color.R == 0 {
		p.color = planetColors[rand.Intn(len(planetColors))]
	}

	return p
}

func (g *Game) GeneratePlanets() {
	for i := 0; i < randomPlanetCount; i++ {
		newPlanet := NewPlanet(
			Planet{
				radius:  rand.Float64() * maxPlanetSize,
				density: planetDensity,

				pos: Pair{
					rand.Float64() * screenWidth,
					rand.Float64() * screenHeight,
				},
			},
		)

		g.state = append(g.state, newPlanet)
	}
}

func (g *Game) GenerateSun() {
	newPlanet := NewPlanet(
		Planet{
			radius:  sunSize,
			density: sunDensity,
			isSun:   true,

			color: planetColors[0],

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

func (g *Game) UpdateAllPlanetPositions() {
	var wg sync.WaitGroup
	for i := range g.state {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			g.calculateNextPlanetPosition(i)
			if g.state[i].collisionCooldown > 0 {
				g.state[i].collisionCooldown--
			}
		}(i)
	}
	wg.Wait()

	for i := range g.state {
		g.state[i].pos = g.state[i].nextPos
	}
}

func (g *Game) calculateNextPlanetPosition(index int) {
	mainPlanet := &g.state[index]
	if mainPlanet.isDead {
		return
	}

	totalForceX, totalForceY, collidedWithSun := calculateForceFromTree(g.root, mainPlanet)
	if collidedWithSun {
		mainPlanet.isDead = true
		return
	}

	mainPlanet.vel.X += totalForceX / mainPlanet.mass
	mainPlanet.vel.Y += totalForceY / mainPlanet.mass

	mainPlanet.nextPos.X = mainPlanet.pos.X + mainPlanet.vel.X
	mainPlanet.nextPos.Y = mainPlanet.pos.Y + mainPlanet.vel.Y
}

func (g *Game) RemoveDeadPlanets() {
	for i, planet := range g.state {
		if planet.isDead && !planet.isSun {
			g.state = append(g.state[:i], g.state[i+1:]...)
		}
	}
}

// func movingTogether(planetA, planetB Planet) bool {
// 	// Relative velocity
// 	dvx := planetB.vel.X - planetA.vel.X
// 	dvy := planetB.vel.Y - planetA.vel.Y

// 	// Relative position
// 	dx := planetB.pos.X - planetA.pos.X
// 	dy := planetB.pos.Y - planetA.pos.Y

// 	// Dot product to check if moving toward each other
// 	dot := dvx*dx + dvy*dy

// 	if dot < 0 {
// 		return false
// 	}

// 	return true
// }

// n^2 algo
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
