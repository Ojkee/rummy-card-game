package window

import rl "github.com/gen2brain/raylib-go/raylib"

type FuncButton struct {
	rect    rl.Rectangle
	content string
	color   rl.Color
}

func NewFuncButton(rect rl.Rectangle, content string) *FuncButton {
	return &FuncButton{
		rect:    rect,
		content: content,
		color:   COLOR_BUTTON_NOT_READY,
	}
}

func (fb *FuncButton) InRect(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, fb.rect)
}

func (fb *FuncButton) UpdateRect(rect *rl.Rectangle) {
	fb.rect = *rect
}

func (fb *FuncButton) UpdateContent(newContent string) {
	fb.content = newContent
}
