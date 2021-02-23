package main


import (
	"fmt"
	// "time"
	"github.com/veandco/go-sdl2/sdl"
)


const (
	WINDOW_WIDTH  = 1400
	WINDOW_HEIGHT = 800
	SPEED         = 10
)


type Color struct {
	R, G, B byte
}


type Position struct {
	x, y float32
}


type Ball struct {
	Position
	radius   int
	xv       float32
	yv       float32
	color    Color
}


func (ball *Ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x * x + y * y < ball.radius * ball.radius {
				setPixel(int(ball.x) + x, int(ball.y) + y, ball.color, pixels)
			}
		}
	}
}


func (ball *Ball) update(leftPaddle, rightPaddle *Paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	if int(ball.y) - ball.radius < 0 || int(ball.y) + ball.radius > WINDOW_HEIGHT {
		ball.yv = -ball.yv
	}

	if ball.x < 0 || int(ball.x) > WINDOW_WIDTH {
		ball.x = WINDOW_WIDTH  / 2
		ball.y = WINDOW_HEIGHT / 2
	}

	if int(ball.x) < int(leftPaddle.x) + leftPaddle.width {
		colisionLeftPaddle := (
			int(ball.y) > int(leftPaddle.y) - leftPaddle.height / 2 &&
			int(ball.y) < int(leftPaddle.y) + leftPaddle.height / 2)

		if colisionLeftPaddle {
			ball.xv = -ball.xv
		}
	}

	if int(ball.x) > int(rightPaddle.x) - rightPaddle.width {
		colisionRightPaddle := (
			int(ball.y) > int(rightPaddle.y) - rightPaddle.height / 2 &&
    	    int(ball.y) < int(rightPaddle.y) + rightPaddle.height / 2)

		if colisionRightPaddle {
			ball.xv = -ball.xv
		}
	}
}


type Paddle struct {
	Position
	width    int
	height   int
	color    Color
}


func (paddle *Paddle) draw(pixels []byte) {
	start_x := int(paddle.x) - paddle.width  / 2
	start_y := int(paddle.y) - paddle.height / 2

	for y := 0; y < paddle.height; y++ {
		for x := 0; x < paddle.width; x++ {
			setPixel(start_x + x, start_y + y, paddle.color, pixels)
		}
	}
}


func (paddle *Paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= SPEED
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += SPEED
	}
}


func (paddle *Paddle) aiUpdate(ball *Ball) {
	paddle.y = ball.y
}


func setPixel(x, y int, color Color, pixels []byte) {
	index := (y * WINDOW_WIDTH + x) * 4

	if index < len(pixels) - 4 && index >= 0 {
		pixels[index]     = color.R
		pixels[index + 1] = color.G
		pixels[index + 2] = color.B
	}
}


func clearScreen(pixels []byte) {
	for index := range pixels {
		pixels[index] = 0
	}
}


func main() {
	window, err := sdl.CreateWindow(
		"PONG",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH,
		WINDOW_HEIGHT,
		sdl.WINDOW_SHOWN,
	)

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


	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		WINDOW_WIDTH,
		WINDOW_HEIGHT,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer texture.Destroy()


	pixels := make([]byte, WINDOW_WIDTH * WINDOW_HEIGHT * 4)


	player1 := Paddle{Position{10, 100}, 20, 150, Color{255, 255, 255}}
	player2 := Paddle{Position{WINDOW_WIDTH - 10, 100}, 20, 150, Color{255, 255, 255}}
	ball    := Ball{Position{300, 300}, 15, 5, 5, Color{255, 255, 255}}


	keyState := sdl.GetKeyboardState()

	
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clearScreen(pixels)

		player1.update(keyState)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		texture.Update(nil, pixels, WINDOW_WIDTH * 4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(5)
	}
}
