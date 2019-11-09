package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type paddle struct {
	position // this will look for a struct called position and bring its attributes in (you also get all of its receiver functions)
	w        float32
	h        float32
	speed    float32
	score    int
	color    color
}

type aiSettings struct {
	speed float32
}

// lerp (linear interpolation) will get you the number in between 2 values
// [1				10]
// 1 + .5 * (10 - 1) = 5.5

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (paddle *paddle) middleZoneSize() float32 {
	return paddle.h * .3
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	// reason to start with y first:
	// cpus are faster than ram.  So whne you ask from a number from RAM, you'll get the next 64bytes (contiguous - and I think in the cache?)
	// in order, so you can go through the bytes faster in order
	// if you did x first, you would do 0, 3, 6, 1, 4, 7, 2, 5, 8
	// and you would miss the cache often and go slower

	// 0 1 2
	// 3 4 5
	// 6 7 9
	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(position{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 { // 0 if key is not being pressed, 0 if it is being pressed
		paddle.y -= paddle.speed * elapsedTime
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}

	// can't trust that joysticks are always standing still - so we want to set a threshold so that the paddles doesn't just float around
	if math.Abs(float64(controllerAxis)) > 1500 {
		// 32767 is the max value of a float32 (i guess) so this is checking how complete of a joystick movement you have done
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime // max speed * pct of the way your joystick is pressed in and then scale it for elapsedTime for how fast the framerates are going
	}

}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if ball.y > paddle.y {
		paddle.y += paddle.speed * elapsedTime
	}
	if ball.y < paddle.y {
		paddle.y -= paddle.speed * elapsedTime
	}
}
