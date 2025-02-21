package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
)

func (window *Window) handleDiscardButton(mousePos *rl.Vector2) {
	if card := window.getCardIfOneSelected(); card != nil &&
		window.discardButton.isClicked(mousePos) {
		discardAction := cm.NewActionDiscardMessage(window.clientId, card.srcCard)
		window.sendActionCallback(discardAction)
	}
}

func (window *Window) updateDiscardButtonPos() {
	newRect := window.discardButton.rect
	for _, card := range window.playerCards {
		if card.isSelected {
			newRect.X = card.rect.X - float32(DISCARD_BUTTON_WIDTH-CARD_WIDTH)/2
			newRect.Y = card.rect.Y - DISCARD_BUTTON_HEIGHT - 24
			break
		}
	}
	window.discardButton.UpdateRect(&newRect)
}
