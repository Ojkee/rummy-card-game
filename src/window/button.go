package window

import rl "github.com/gen2brain/raylib-go/raylib"

type FuncButton struct {
	rect    rl.Rectangle
	color   rl.Color
	content string
}

func NewFuncButton() *FuncButton {
	rectRect := rl.NewRectangle(
		float32(WINDOW_WIDTH-READY_BUTTON_WIDTH)/2,
		float32(WINDOW_HEIGHT-READY_BUTTON_HEIGHT)/2,
		READY_BUTTON_WIDTH,
		READY_BUTTON_HEIGHT,
	)
	return &FuncButton{
		rect:    rectRect,
		color:   COLOR_BUTTON_NOT_READY,
		content: "Not ready",
	}
}

func (fb *FuncButton) isClicked(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, fb.rect)
}
