package main

import "math"

const (
	// Gravitational Constant
	G = 6.67430e-11
)

func CalculateForcesSlow(mainPlanet, otherPlanet Planet) (xForce float32, yForce float32, collision bool) {
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
	m1 := mainPlanet.mass
	m2 := otherPlanet.mass

	xInelastic := (m1*mainPlanet.vel.X + m2*otherPlanet.vel.X) / (m1 + m2)
	yInelastic := (m1*mainPlanet.vel.Y + m2*otherPlanet.vel.Y) / (m1 + m2)

	xElastic := ((m1-m2)*mainPlanet.vel.X + 2*m2*otherPlanet.vel.X) / (m1 + m2)
	yElastic := ((m1-m2)*mainPlanet.vel.Y + 2*m2*otherPlanet.vel.Y) / (m1 + m2)

	xVel = restitution*xElastic + (1-restitution)*xInelastic
	yVel = restitution*yElastic + (1-restitution)*yInelastic

	return xVel, yVel
}
