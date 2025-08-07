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

	// _, wheelY := ebiten.Wheel()
	// if wheelY > 0 {
	// 	cam.zoom += wheelY * 0.03
	// 	if cam.zoom < 0.1 {
	// 		cam.zoom = 0.1
	// 	}
	// }

	mx, my := ebiten.CursorPosition()
	mouseX := float64(mx)
	mouseY := float64(my)

	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		zoomFactor := 1.0 + wheelY*0.1
		oldZoom := cam.zoom
		newZoom := cam.zoom * zoomFactor

		// Calculate the world position under the mouse before zoom
		worldXBefore := cam.pos.X + mouseX/oldZoom
		worldYBefore := cam.pos.Y + mouseY/oldZoom

		// Apply new zoom
		cam.zoom = newZoom
		if cam.zoom < 0.1 {
			cam.zoom = 0.1
		}

		// Calculate the world position under the mouse after zoom
		worldXAfter := cam.pos.X + mouseX/cam.zoom
		worldYAfter := cam.pos.Y + mouseY/cam.zoom

		// Offset camera to keep mouse-over-world position stable
		cam.pos.X += worldXBefore - worldXAfter
		cam.pos.Y += worldYBefore - worldYAfter
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
