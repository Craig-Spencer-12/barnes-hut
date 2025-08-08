package main

import "image/color"

const (

	// Default Screen
	screenHeight = 1024
	screenWidth  = 2048

	// Initialize Planets
	randomPlanetCount = 10000
	maxPlanetSize     = 10
	planetDensity     = 16000
	sunSize           = 70
	sunDensity        = 2000000000
	startingVelocity  = 5

	// Default Start State
	startPaused    = false
	startDrawLines = false

	// Collision Tweaks
	collisionCooldown = 5
	collisionsOn      = false

	// Math
	G           = 6.67430e-11 // Gravitational Constant
	theta       = 0.5         // tuning parameter for accuracy vs performance
	restitution = 0.6         // how elastic are collisions
	minDistance = 0.0001      // prevents infinite recursion when building the tree
	nearZero    = 1e-4
)

var sunColor = color.RGBA{R: 255, G: 215, B: 0, A: 255}
var planetColors = [...]color.RGBA{
	{R: 255, G: 215, B: 0, A: 255},   // Gold (Venus-like)
	{R: 70, G: 130, B: 180, A: 255},  // Steel blue (Earth-like)
	{R: 205, G: 92, B: 92, A: 255},   // Indian red (Mars-like)
	{R: 100, G: 149, B: 237, A: 255}, // Cornflower blue (Neptune)
	{R: 210, G: 180, B: 140, A: 255}, // Tan (rocky)
	{R: 128, G: 0, B: 128, A: 255},   // Purple (exotic)
}
