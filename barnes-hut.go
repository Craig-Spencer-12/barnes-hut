package main

import (
	"math"
)

func (g *Game) BuldTree() {
	topLeft, size := g.computeBounds()
	g.root = &Node{
		topLeft: topLeft,
		size:    size,
	}

	for i := range g.state {
		g.root.insert(&g.state[i])
	}
}

func CalculateCentersOfMassRecursive(n *Node) {
	if n == nil {
		return
	}

	if n.body != nil {
		n.totalMass = n.body.mass
		n.centerOfMass = n.body.pos
		return
	}

	var numeratorX, numeratorY float64 = 0, 0
	var denominator float64 = 0
	for _, node := range n.child {
		if node != nil {
			if n.totalMass == 0 {
				CalculateCentersOfMassRecursive(node)
			}
			numeratorX += node.totalMass * node.centerOfMass.X
			numeratorY += node.totalMass * node.centerOfMass.Y
			denominator += node.totalMass
		}
	}

	n.centerOfMass.X = numeratorX / denominator
	n.centerOfMass.Y = numeratorY / denominator
	n.totalMass = denominator
}

func (n *Node) insert(planet *Planet) {

	// Find which region planet is in
	region := 0
	newTopLeft := n.topLeft
	newSize := n.size / 2
	if planet.pos.X > n.topLeft.X+newSize {
		region++
		newTopLeft.X += newSize
	}
	if planet.pos.Y > n.topLeft.Y+newSize {
		region += 2
		newTopLeft.Y += newSize
	}

	// Add leaf if theres room
	if n.child[region] == nil {
		n.child[region] = &Node{
			parent: n,
			body:   planet,

			topLeft: newTopLeft,
			size:    newSize,
		}

		return
	}

	// If desired leaf is full, subdivide region
	// replace the node with a subdivided node and insert both nodes on that  subdivided node
	existingPlanet := n.child[region].body
	if existingPlanet != nil {
		n.child[region] = &Node{
			parent:  n,
			topLeft: newTopLeft,
			size:    newSize,
		}

		// Edge Case: planets too close together cause infinte recursion
		dx := math.Abs(existingPlanet.pos.X - planet.pos.X)
		dy := math.Abs(existingPlanet.pos.Y - planet.pos.Y)
		if dx > minDistance || dy > minDistance {
			n.child[region].insert(existingPlanet)
		}
	}

	n.child[region].insert(planet)
}

func (g *Game) computeBounds() (topLeft Pair, size float64) {
	if len(g.state) == 0 {
		return Pair{0, 0}, 0
	}

	minX, maxX := g.state[0].pos.X, g.state[0].pos.X
	minY, maxY := g.state[0].pos.Y, g.state[0].pos.Y

	for _, p := range g.state {
		if p.pos.X < minX {
			minX = p.pos.X
		}
		if p.pos.X > maxX {
			maxX = p.pos.X
		}
		if p.pos.Y < minY {
			minY = p.pos.Y
		}
		if p.pos.Y > maxY {
			maxY = p.pos.Y
		}
	}

	width := maxX - minX
	height := maxY - minY
	size = math.Max(width, height)
	buffer := size * 0.1

	return Pair{minX - buffer, minY - buffer}, size + buffer*2
}

func calculateForceFromTree(n *Node, target *Planet) (fx, fy float64, collided bool) {
	if n == nil || (n.body == target && n.body != nil) {
		return 0, 0, false
	}

	// Compute distance from center of mass to target
	dx := n.centerOfMass.X - target.pos.X
	dy := n.centerOfMass.Y - target.pos.Y
	distSq := dx*dx + dy*dy
	dist := math.Sqrt(distSq)

	// Collision detection
	if n.body != nil && dist < target.radius+n.body.radius {
		if n.body.isSun {
			return 0, 0, true
		}

		// if target.collisionCooldown == 0 && !movingTogether(*n.body, *target) {
		// 	vx, vy := collisionVelocities(*target, *n.body)
		// 	target.vel.X = vx
		// 	target.vel.Y = vy
		// 	target.collisionCooldown = collisionCooldown
		// 	return 0, 0, false
		// }
	}

	if n.totalMass == 0 {
		return 0, 0, false
	}

	if n.body != nil || (n.size/dist < theta) {
		fx, fy := CalculateGravityForce(*target, n.centerOfMass, n.totalMass)
		return fx, fy, false
	}

	var totalX, totalY float64
	for _, child := range n.child {
		cx, cy, ccol := calculateForceFromTree(child, target)
		totalX += cx
		totalY += cy
		if ccol {
			return 0, 0, true
		}
	}
	return totalX, totalY, false
}

func CalculateGravityForce(target Planet, otherPos Pair, otherMass float64) (fx, fy float64) {

	dx := otherPos.X - target.pos.X
	dy := otherPos.Y - target.pos.Y
	distSq := dx*dx + dy*dy + nearZero
	dist := math.Sqrt(distSq)

	force := G * target.mass * otherMass / distSq
	fx = force * dx / dist
	fy = force * dy / dist
	return
}
