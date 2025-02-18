package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"rummy-card-game/src/connection_messages"
)

func (window *Window) inRoundManager(mousePos *rl.Vector2) {
	if window.drawPile.IsClicked(mousePos) {
		actionMsg := connection_messages.NewActionDrawMessage(window.clientId)
		window.sendActionCallback(actionMsg)
	}
}

func (window *Window) drawInRound() {
	for _, playerCard := range window.playerCards {
		playerCard.Draw()
	}
	window.lastDiscardedCard.Draw()
	window.drawPile.Draw()
	window.drawDisplayText()
}

func (window *Window) drawDisplayText() {
	if window.displayTime > 0 {
		rl.DrawTextEx(
			FONT,
			window.displayText,
			rl.NewVector2(10, 10),
			float32(FONT_SIZE),
			FONT_SPACING,
			COLOR_BEIGE,
		)
		window.displayTime -= rl.GetFrameTime()
	}
}
