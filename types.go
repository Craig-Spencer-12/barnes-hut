package main

import (
	"image/color"
)

type Game struct {
	state []Planet
	root  *Node

	camera    Camera
	paused    bool
	drawLines bool
	stepFrame bool
}

type Planet struct {
	radius            float64
	density           float64
	mass              float64
	isSun             bool
	isDead            bool
	collisionCooldown int

	vel     Pair
	pos     Pair
	nextPos Pair

	color color.RGBA
}

type Pair struct {
	X float64
	Y float64
}

type Node struct {
	parent *Node
	child  [4]*Node

	topLeft Pair
	size    float64

	body *Planet

	centerOfMass Pair
	totalMass    float64
}
