package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	width       int32
	height      int32
	playerCards []CardModel
	discardCard CardModel
}

func NewWindow() *Window {
	return &Window{
		width:  800,
		height: 800,
	}
}

func (window *Window) MainLoop() {
	rl.InitWindow(window.width, window.height, "Rummy")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		window.checkEvent()
		window.draw()
	}
}

func (window *Window) checkEvent() {
}

func (window *Window) draw() {
	rl.BeginDrawing()

	rl.ClearBackground(COLOR_DARK_GRAY)

	rl.EndDrawing()
}
