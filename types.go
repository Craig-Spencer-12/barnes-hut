package main

type Planet struct {
	radius  float32
	density float32
	mass    float32

	vel     Pair
	pos     Pair
	nextPos Pair
}

type Pair struct {
	X float32
	Y float32
}
