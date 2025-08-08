package main

import "math"

func CalculateForcesSlow(mainPlanet, otherPlanet Planet) (xForce float64, yForce float64, collision bool) {
	dx := otherPlanet.pos.X - mainPlanet.pos.X
	dy := otherPlanet.pos.Y - mainPlanet.pos.Y
	distance := math.Hypot(dx, dy)

	// Collision Check
	if distance < mainPlanet.radius+otherPlanet.radius {
		return 0, 0, true
	}

	force := G * mainPlanet.mass * otherPlanet.mass / (distance * distance)

	xForce = force * dx / distance
	yForce = force * dy / distance

	return xForce, yForce, false
}

func collisionVelocities(mainPlanet, otherPlanet Planet) (xVel, yVel float64) {
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
