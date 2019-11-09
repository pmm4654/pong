package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// frame rate independence
// score
// game over state - win/lose
// AI needs to be able to lose (more imperfect!)
// Handling resizing of the window

const winWidth int = 800
const winHeight int = 600

type color struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

// -- this is the syntax for emum
type gameState int

const (
	start gameState = iota // 0
	play                   // 1
)

//

var state = start

func setPixel(x int, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	// index would be whateveryou y is, so say 1 and your x is 2, you would
	// want multiply by your width (3 in the example below) and add the x (2).
	// So you would multiply 1 * 3 + 2 and your pixel would be placed at 5 (x = 2,y = 1)
	//
	// 0 1 2
	// 3 4 [5]
	// 6 7 8
	//
	// Multiplying by 4 is because there are 4 bytes per pixel in the format sdl.PIXELFORMAT_ABGR8888
	// The 4 bytes are taken up by alpha, red, greed and blue

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer texture.Destroy()

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, winWidth*winHeight*4)

	// Big Game Loop
	currentAi := aiSettings{speed: 300}

	player1 := paddle{position{50, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	player2 := paddle{position{float32(winWidth) - 50, 100}, 20, 100, currentAi.speed, 0, color{255, 255, 255}}
	ball := ball{position{300, 300}, 20, 400, 400, color{255, 255, 255}}

	// this is balically an array that has arepresentation of every key and whether or not a key is being pressed
	// it keeps you from having to check each event manually yourslef in the game loop
	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32
	var controllerAxis int16

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// for joysticks
		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		if state == play {
			drawNumber(getCenter(), color{255, 255, 255}, 20, player1.score, pixels)
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 { // reset score
				if player1.score == 3 || player2.score == 3 { // only reset the score if one player has won (3 is winning here)
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 { // if elapsed time was 5 ms, about 200 fps (1 / .005)
			sdl.Delay(5 - uint32(elapsedTime/1000.0)) // delay until we get to 5 ms
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

	// sdl.Delay(2000)
}

// game loops

// for {
// Commonly broken up into update and draw
// Update
// get all input (keyboard/mouse etc)
// update all yoru things (physics, ai, whatever)
// Draw
// draw all the things
// }
