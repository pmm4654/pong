package main

import "github.com/veandco/go-sdl2/sdl"

type paddle struct {
	position // this will look for a struct called position and bring its attributes in (you also get all of its receiver functions)
	w        int
	h        int
	color    color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - int(paddle.w)/2
	startY := int(paddle.y) - int(paddle.h)/2

	// reason to start with y first:
	// cpus are faster than ram.  So whne you ask from a number from RAM, you'll get the next 64bytes (contiguous - and I think in the cache?)
	// in order, so you can go through the bytes faster in order
	// if you did x first, you would do 0, 3, 6, 1, 4, 7, 2, 5, 8
	// and you would miss the cache often and go slower

	// 0 1 2
	// 3 4 5
	// 6 7 9
	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 { // 0 if key is not being pressed, 0 if it is being pressed
		paddle.y -= 5
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 5
	}
}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}
