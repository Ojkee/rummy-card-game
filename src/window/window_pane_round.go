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
}
