package main

import "math"

const (
	// Gravitational Constant
	G = 6.67430e-11
)

func CalculateForces(mainPlanet, otherPlanet Planet) (xForce float32, yForce float32, collision bool) {
	dx := float64(otherPlanet.pos.X - mainPlanet.pos.X)
	dy := float64(otherPlanet.pos.Y - mainPlanet.pos.Y)
	distance := float32(math.Hypot(dx, dy))

	// Collision Check
	if distance < mainPlanet.radius+otherPlanet.radius {
		return 0, 0, true
	}

	force := G * mainPlanet.mass * otherPlanet.mass / (distance * distance)

	xForce = force * float32(dx) / distance
	yForce = force * float32(dy) / distance

	return xForce, yForce, false
}

func collisionVelocities(mainPlanet, otherPlanet Planet) (xVel, yVel float32) {
	// Elastic
	// xVel = ((mainPlanet.mass-otherPlanet.mass)*mainPlanet.vel.X + (2 * otherPlanet.mass * otherPlanet.vel.X)) / (mainPlanet.mass + otherPlanet.mass)
	// yVel = ((mainPlanet.mass-otherPlanet.mass)*mainPlanet.vel.Y + (2 * otherPlanet.mass * otherPlanet.vel.Y)) / (mainPlanet.mass + otherPlanet.mass)

	// Inelastic
	xVel = (mainPlanet.mass*mainPlanet.vel.X + otherPlanet.mass*otherPlanet.vel.X) / (mainPlanet.mass + otherPlanet.mass)
	yVel = (mainPlanet.mass*mainPlanet.vel.Y + otherPlanet.mass*otherPlanet.vel.Y) / (mainPlanet.mass + otherPlanet.mass)

	return xVel, yVel
}
