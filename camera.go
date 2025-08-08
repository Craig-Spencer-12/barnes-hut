package main

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	pos  Pair
	zoom float64

	isDragging     bool
	prevMousePoint Pair
}

func (c *Camera) Controls() {
	c.wasdMove()
	c.scrollZoom()
	c.dragCamera()
}

func (c *Camera) wasdMove() {
	speed := 10 / c.zoom

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.pos.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.pos.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.pos.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.pos.X += speed
	}
}

func (c *Camera) scrollZoom() {
	mx, my := ebiten.CursorPosition()
	mouseX := float64(mx)
	mouseY := float64(my)

	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {

		worldXBefore := c.pos.X + mouseX/c.zoom
		worldYBefore := c.pos.Y + mouseY/c.zoom

		c.zoom += wheelY * 0.1
		if c.zoom < 0.1 {
			c.zoom = 0.1
		}

		worldXAfter := c.pos.X + mouseX/c.zoom
		worldYAfter := c.pos.Y + mouseY/c.zoom

		c.pos.X += worldXBefore - worldXAfter
		c.pos.Y += worldYBefore - worldYAfter
	}
}

func (c *Camera) dragCamera() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		mouse := Pair{float64(x), float64(y)}

		if !c.isDragging {
			c.isDragging = true
			c.prevMousePoint = mouse
		} else {
			dx := (mouse.X - c.prevMousePoint.X) / c.zoom
			dy := (mouse.Y - c.prevMousePoint.Y) / c.zoom
			c.pos.X -= dx
			c.pos.Y -= dy
			c.prevMousePoint = mouse
		}
	} else {
		c.isDragging = false
	}
}
