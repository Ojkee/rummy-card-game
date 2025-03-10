package window

import rl "github.com/gen2brain/raylib-go/raylib"

type DrawPileButton struct {
	rect rl.Rectangle
}

func NewDrawPileButton() *DrawPileButton {
	return &DrawPileButton{
		rect: DRAW_PILE_POS,
	}
}

func (drawPile *DrawPileButton) InRect(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, drawPile.rect)
}

func (drawPile *DrawPileButton) Draw() {
	rl.DrawRectangleRec(
		drawPile.rect,
		COLOR_BEIGE,
	)
	rl.DrawRectangleRec(
		rl.NewRectangle(
			drawPile.rect.X+float32(CARD_GAP),
			float32(drawPile.rect.ToInt32().Y+CARD_GAP),
			float32(CARD_INNER_WIDTH),
			float32(CARD_INNER_HEIGHT),
		),
		COLOR_TAUPE,
	)
	drawTextSize := GetTextVec("Draw")
	rl.DrawTextEx(
		FONT,
		"Draw",
		rl.NewVector2(
			drawPile.rect.X+(drawPile.rect.Width-drawTextSize.X)/2,
			drawPile.rect.Y-drawTextSize.Y,
		),
		FONT_SIZE,
		FONT_SPACING,
		COLOR_BEIGE,
	)
}
