package main

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	pos  Pair
	zoom float64

	isDragging      bool
	prevMouseScreen Pair
}

func (cam *Camera) Update() {
	speed := 10 / cam.zoom

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cam.pos.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cam.pos.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cam.pos.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cam.pos.X += speed
	}

	_, wheelY := ebiten.Wheel()
	cam.zoom += float64(wheelY) * 0.03
	if cam.zoom < 0.1 {
		cam.zoom = 0.1
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		mouse := Pair{float64(x), float64(y)}

		if !cam.isDragging {
			cam.isDragging = true
			cam.prevMouseScreen = mouse
		} else {
			dx := (mouse.X - cam.prevMouseScreen.X) / cam.zoom
			dy := (mouse.Y - cam.prevMouseScreen.Y) / cam.zoom
			cam.pos.X -= dx
			cam.pos.Y -= dy
			cam.prevMouseScreen = mouse
		}
	} else {
		cam.isDragging = false
	}
}
