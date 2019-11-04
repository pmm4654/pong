package main

type ball struct {
	position
	radius int
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
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle) {
	ball.x += ball.xVel
	ball.y += ball.yVel

	// handle colisions with top and bottom (bounce it off)
	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > winHeight {
		ball.yVel = -ball.yVel
	}

	if int(ball.x)-ball.radius < 0 || int(ball.x)+ball.radius > winWidth {
		ball.x = 300
		ball.y = 300
	}

	if int(ball.x) < int(leftPaddle.x)+leftPaddle.w/2 {
		if int(ball.y) > int(leftPaddle.y)-leftPaddle.h/2 && int(ball.y) < int(leftPaddle.y)+leftPaddle.h/2 {
			ball.xVel = -ball.xVel
		}
	}

	if int(ball.x) > int(rightPaddle.x)-rightPaddle.w/2 {
		if int(ball.y) > int(rightPaddle.y)-rightPaddle.h/2 && int(ball.y) < int(rightPaddle.y)+rightPaddle.h/2 {
			ball.xVel = -ball.xVel
		}
	}
}
