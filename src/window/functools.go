package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GetTextVec(text string) rl.Vector2 {
	return rl.MeasureTextEx(FONT, text, FONT_SIZE, FONT_SPACING)
}
