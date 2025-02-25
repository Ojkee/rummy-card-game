package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
)

func (window *Window) inRoundManagerClick(mousePos *rl.Vector2) {
	if window.drawPile.InRect(mousePos) {
		actionMsg := cm.NewActionDrawMessage(
			window.clientId,
			cm.DRAW_FROM_PILE,
		)
		window.sendActionCallback(actionMsg)
	} else if window.lastDiscardedCard.srcCard != nil &&
		window.lastDiscardedCard.InRect(*mousePos) {
		actionMsg := cm.NewActionDrawMessage(
			window.clientId,
			cm.DRAW_FROM_DISCARD_PILE,
		)
		window.sendActionCallback(actionMsg)
	}
	window.handleCardClicked(mousePos)
	window.handleDiscardButton(mousePos)
	window.handleLockSequence(mousePos)
	window.handleInitialMeldButton(mousePos)
	window.handleAvailableSpots(mousePos)
}

func (window *Window) drawInRound() {
	for i, playerCard := range window.playerCards {
		if i != window.currentDragCardIdx {
			playerCard.Draw()
		}
	}
	if window.currentDragCardIdx != -1 {
		window.playerCards[window.currentDragCardIdx].Draw()
	}
	window.lastDiscardedCard.Draw()
	window.drawPile.Draw()
	window.drawDisplayText()
	window.drawTurnInfo()
	if card := window.getCardIfOneSelected(); card != nil {
		window.updateDiscardButtonPos()
		window.drawStaticButton(&window.discardButton)
	}
	if selectedCards := window.getSelectedUnlockedCards(); len(selectedCards) >= 3 {
		window.drawStaticButton(&window.lockSetButton)
	}
	if window.numLockedSequences() > 0 {
		window.drawStaticButton(&window.initialMeldButton)
	}
	for _, sequence := range window.tableSequences {
		sequence.Draw()
	}
	for _, availableSpot := range window.availableSpots {
		availableSpot.Draw()
	}
}

func (window *Window) handleCardClicked(mousePos *rl.Vector2) {
	for i := range window.playerCards {
		if window.playerCards[i].InRect(*mousePos) {
			window.playerCards[i].isSelected = !window.playerCards[i].isSelected
			if window.playerCards[i].sequenceId != -1 {
				window.unlockAllById(window.playerCards[i].sequenceId)
			}
			if card := window.getCardIfOneSelected(); card != nil {
				window.initAvailableSpots(card)
			} else {
				window.resetAvailableSpots()
			}
			break
		}
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
