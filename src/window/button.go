package window

import rl "github.com/gen2brain/raylib-go/raylib"

type FuncButton struct {
	rect    rl.Rectangle
	color   rl.Color
	content string
}

func NewFuncButton(rect rl.Rectangle, content string) *FuncButton {
	return &FuncButton{
		rect:    rect,
		color:   COLOR_BUTTON_NOT_READY,
		content: content,
	}
}

func (fb *FuncButton) isClicked(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, fb.rect)
}

func (fb *FuncButton) UpdateRect(rect *rl.Rectangle) {
	fb.rect = *rect
}
