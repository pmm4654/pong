package main

import (
	"fmt"

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

	pixels := make([]byte, winWidth*winHeight*4)

	// Big Game Loop

	player1 := paddle{position{50, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{position{float32(winWidth) - 50, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{position{300, 300}, 20, 2, 2, color{255, 255, 255}}

	// this is balically an array that has arepresentation of every key and whether or not a key is being pressed
	// it keeps you from having to check each event manually yourslef in the game loop
	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		player1.update(keyState)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(16)
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
