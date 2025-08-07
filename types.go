package main

import "image/color"

type Planet struct {
	radius            float32
	density           float32
	mass              float32
	isSun             bool
	isDead            bool
	collisionCooldown int

	vel     Pair
	pos     Pair
	nextPos Pair

	color color.RGBA
}

type Pair struct {
	X float32
	Y float32
}

var planetColors = [6]color.RGBA{
	{R: 70, G: 130, B: 180, A: 255},  // Steel blue (Earth-like)
	{R: 205, G: 92, B: 92, A: 255},   // Indian red (Mars-like)
	{R: 255, G: 215, B: 0, A: 255},   // Gold (Venus/gas giant)
	{R: 100, G: 149, B: 237, A: 255}, // Cornflower blue (Neptune)
	{R: 210, G: 180, B: 140, A: 255}, // Tan (rocky)
	{R: 128, G: 0, B: 128, A: 255},   // Purple (exotic)
}
