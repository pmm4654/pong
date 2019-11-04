package main

type ball struct {
	position
	radius float32
	xVel   float32 // velocity
	yVel   float32 // velocity
	color  color
}

func (ball *ball) draw(pixels []byte) {
	// YAGNI - Ya ain't gonna need it - don't pre-maturely optimize
	// Draw this entire rectangle and if the pixels is outside of the radius of the start point, don't draw it
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius { // if x is outside of the radius
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func getCenter() position {
	return position{float32(winWidth) / 2, float32(winHeight) / 2}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xVel * elapsedTime
	ball.y += ball.yVel * elapsedTime

	// handle colisions with top and bottom (bounce it off)
	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(winHeight) {
		ball.yVel = -ball.yVel
	}

	if ball.x-ball.radius < 0 || ball.x+ball.radius > float32(winWidth) {
		ball.position = getCenter()
		// ball.x = 300
		// ball.y = 300
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.xVel = -ball.xVel
		}
	}

	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xVel = -ball.xVel
		}
	}
}
