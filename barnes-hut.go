package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Node struct {
	parent *Node
	child  [4]*Node

	topLeft Pair
	size    float32

	body *Planet

	centerOfMass Pair
	totalMass    float32
}

func (g *Game) CreateTree() {
	topLeft, size := g.computeBounds()
	g.root = &Node{
		topLeft: topLeft,
		size:    size,
	}

	for i := range g.state {
		g.root.insert(&g.state[i])
	}
}

func CalculateCentersOfMass(n *Node) {
	if n == nil {
		return
	}

	if n.body != nil {
		n.totalMass = n.body.mass
		n.centerOfMass = n.body.pos
		return
	}

	var numeratorX, numeratorY float32 = 0, 0
	var denominator float32 = 0
	for _, node := range n.child {
		if node != nil {
			if n.totalMass == 0 {
				CalculateCentersOfMass(node)
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
	// which quad is new node
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

	// if that leaf is empty just place the planet there
	if n.child[region] == nil {
		n.child[region] = &Node{
			parent: n,
			body:   planet,

			topLeft: newTopLeft,
			size:    newSize,
		}

		return
	}

	// else if that leaf is full we need to subdivide that region until they occupy different regions
	// replace the node with a subdivided node and insert both nodes on that  subdivided node

	existingPlanet := n.child[region].body
	if existingPlanet != nil {
		n.child[region] = &Node{
			parent:  n,
			topLeft: newTopLeft,
			size:    newSize,
		}

		// making sure planets aren't on top of each other causing infinte recursion
		dx := math.Abs(float64(existingPlanet.pos.X - planet.pos.X))
		dy := math.Abs(float64(existingPlanet.pos.Y - planet.pos.Y))
		if dx > minDistance || dy > minDistance {
			n.child[region].insert(existingPlanet)
		}
	}

	n.child[region].insert(planet)
}

func (g *Game) computeBounds() (topLeft Pair, size float32) {
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
	size = float32(math.Max(float64(width), float64(height)))
	buffer := size * 0.1

	return Pair{minX - buffer, minY - buffer}, size + buffer*2
}

func drawBarnesHutBorders(screen *ebiten.Image, n *Node) {
	if n == nil {
		return
	}

	vector.StrokeLine(screen, n.topLeft.X, n.topLeft.Y, n.topLeft.X+n.size, n.topLeft.Y, 1, color.White, true)
	vector.StrokeLine(screen, n.topLeft.X, n.topLeft.Y, n.topLeft.X, n.topLeft.Y+n.size, 1, color.White, true)
	vector.StrokeLine(screen, n.topLeft.X+n.size, n.topLeft.Y+n.size, n.topLeft.X, n.topLeft.Y+n.size, 1, color.White, true)
	vector.StrokeLine(screen, n.topLeft.X+n.size, n.topLeft.Y+n.size, n.topLeft.X+n.size, n.topLeft.Y, 1, color.White, true)

	for i := range n.child {
		drawBarnesHutBorders(screen, n.child[i])
	}
}

func (g *Game) updatePlanetPosition(index int) {
	mainPlanet := &g.state[index]
	if mainPlanet.isDead {
		return
	}

	totalForceX, totalForceY, collidedWithSun := calculateForceFromTree(g.root, mainPlanet, theta)
	if collidedWithSun {
		mainPlanet.isDead = true
		return
	}

	mainPlanet.vel.X += totalForceX / mainPlanet.mass
	mainPlanet.vel.Y += totalForceY / mainPlanet.mass

	mainPlanet.nextPos.X = mainPlanet.pos.X + mainPlanet.vel.X
	mainPlanet.nextPos.Y = mainPlanet.pos.Y + mainPlanet.vel.Y
}

func calculateForceFromTree(n *Node, target *Planet, θ float32) (fx, fy float32, collided bool) {
	if n == nil || (n.body == target && n.body != nil) {
		return 0, 0, false
	}

	// Compute distance from center of mass to target
	dx := n.centerOfMass.X - target.pos.X
	dy := n.centerOfMass.Y - target.pos.Y
	distSq := dx*dx + dy*dy
	dist := float32(math.Sqrt(float64(distSq)))

	// Collision detection (optional: tune threshold)
	if n.body != nil && dist < target.radius+n.body.radius {
		if n.body.isSun {
			return 0, 0, true
		}

		if target.collisionCooldown == 0 && !movingTogether(*n.body, *target) {
			vx, vy := collisionVelocities(*target, *n.body)
			target.vel.X = vx
			target.vel.Y = vy
			target.collisionCooldown = collisionCooldown
			return 0, 0, false
		}
	}

	// Empty node (no mass)
	if n.totalMass == 0 {
		return 0, 0, false
	}

	// Is this region far enough to approximate?
	if n.body != nil || (n.size/dist < θ) {
		// Approximate force from this node's center of mass
		fx, fy := CalculateGravityForce(*target, n.centerOfMass, n.totalMass)
		return fx, fy, false
	}

	// Otherwise, recursively calculate from children
	var totalX, totalY float32
	for _, child := range n.child {
		cx, cy, ccol := calculateForceFromTree(child, target, θ)
		totalX += cx
		totalY += cy
		if ccol {
			return 0, 0, true
		}
	}
	return totalX, totalY, false
}

func CalculateGravityForce(target Planet, otherPos Pair, otherMass float32) (fx, fy float32) {
	const G = 6.674e-11 // or scale to your simulation

	dx := otherPos.X - target.pos.X
	dy := otherPos.Y - target.pos.Y
	distSq := dx*dx + dy*dy + 1e-4 // avoid divide-by-zero
	dist := float32(math.Sqrt(float64(distSq)))

	force := G * target.mass * otherMass / distSq
	fx = force * dx / dist
	fy = force * dy / dist
	return
}
