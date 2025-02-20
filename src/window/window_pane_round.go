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
	window.handleCardClicked(mousePos)
	window.handleDiscardButton(mousePos)
}

func (window *Window) handleCardClicked(mousePos *rl.Vector2) {
	for i := range window.playerCards {
		if window.playerCards[i].IsClicked(*mousePos) {
			window.playerCards[i].isSelected = !window.playerCards[i].isSelected
		}
	}
}

func (window *Window) handleDiscardButton(mousePos *rl.Vector2) {
	if card := window.getCardIfOneSelected(); card != nil &&
		window.discardButton.isClicked(mousePos) {
		discardAction := connection_messages.NewActionDiscardMessage(window.clientId, card.srcCard)
		window.sendActionCallback(discardAction)
	}
}

func (window *Window) getCardIfOneSelected() *CardModel {
	selectedIdx := -1
	for i, card := range window.playerCards {
		if card.isSelected && selectedIdx != -1 {
			return nil
		} else if card.isSelected {
			selectedIdx = i
		}
	}
	if selectedIdx == -1 {
		return nil
	}
	return &window.playerCards[selectedIdx]
}

func (window *Window) drawInRound() {
	for _, playerCard := range window.playerCards {
		playerCard.Draw()
	}
	window.lastDiscardedCard.Draw()
	window.drawPile.Draw()
	window.drawDisplayText()
	window.drawTurnInfo()
	if card := window.getCardIfOneSelected(); card != nil {
		window.drawDiscardButton()
	}
}

func (window *Window) drawDisplayText() {
	if window.displayTime > 0 {
		rl.DrawTextEx(
			FONT,
			window.displayText,
			rl.NewVector2(10, 30),
			float32(FONT_SIZE),
			FONT_SPACING,
			COLOR_BEIGE,
		)
		window.displayTime -= rl.GetFrameTime()
	}
}

func (window *Window) drawTurnInfo() {
	if window.currentTurnId == window.clientId {
		rl.DrawTextEx(
			FONT,
			"Your turn",
			rl.NewVector2(10, 10),
			float32(FONT_SIZE),
			FONT_SPACING,
			COLOR_BEIGE,
		)
	}
}

func (window *Window) drawDiscardButton() {
	var rectOuter rl.Rectangle
	for _, card := range window.playerCards {
		if card.isSelected {
			rectOuter.X = card.rect.X - float32(DISCARD_BUTTON_WIDTH-CARD_WIDTH)/2
			rectOuter.Y = card.rect.Y - DISCARD_BUTTON_HEIGHT - 24
		}
	}
	rectOuter.Width = DISCARD_BUTTON_WIDTH
	rectOuter.Height = DISCARD_BUTTON_HEIGHT
	rectInner := rectOuter
	rectInner.X += 2
	rectInner.Y += 2
	rectInner.Width -= 4
	rectInner.Height -= 4
	rl.DrawRectangleRounded(rectOuter, 0.5, 10, COLOR_BEIGE)
	rl.DrawRectangleRounded(rectInner, 0.5, 10, COLOR_DARK_GRAY)
	rl.DrawTextEx(
		FONT,
		window.discardButton.content,
		rl.NewVector2(rectOuter.X+4, rectOuter.Y+2),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
	window.discardButton.UpdateRect(&rectOuter)
}
