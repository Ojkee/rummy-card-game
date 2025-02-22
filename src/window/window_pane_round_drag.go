package window

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
)

func (window *Window) inRoundManagerDrag(mousePos *rl.Vector2) {
	if window.currentDragCardIdx == -1 {
		window.currentDragCardIdx = window.getDragCardIdx(mousePos)
	} else {
		window.playerCards[window.currentDragCardIdx].MoveX(mousePos.X - float32(CARD_WIDTH)/2)
	}
}

func (window *Window) getDragCardIdx(mousePos *rl.Vector2) int {
	for i, card := range window.playerCards {
		if card.InRect(*mousePos) {
			return i
		}
	}
	return -1
}

func (window *Window) rearrangeNewCardPosX() {
	if window.currentDragCardIdx == -1 {
		return
	}
	newIdx := window.getReleaseDragCardIdx()
	log.Println(newIdx)
	window.insertDragCardNewIdx(newIdx)
	window.sendRearrangedHand()
}

func (window *Window) getReleaseDragCardIdx() int {
	x := window.playerCards[window.currentDragCardIdx].rect.X
	numCards := len(window.playerCards)
	left := float32(WINDOW_WIDTH-int32(numCards)*CARD_WIDTH) / 2
	right := float32(WINDOW_WIDTH+int32(numCards)*CARD_WIDTH) / 2
	if x >= right {
		return numCards
	}
	j := 0
	for i := left; i < right; i += float32(CARD_WIDTH) {
		if x < i {
			return j
		}
		j++
	}
	return numCards - 1
}

func (window *Window) insertDragCardNewIdx(newIdx int) {
	oldIdx := window.currentDragCardIdx
	if oldIdx == newIdx {
		return
	}
	if newIdx == len(window.playerCards) {
		window.playerCards = window.appendNewFilterOld(oldIdx)
	} else {
		window.playerCards = window.insertBeforeIdx(oldIdx, newIdx)
	}
}

func (window *Window) appendNewFilterOld(oldIdx int) []CardModel {
	newHand := make([]CardModel, 0)
	for i, card := range window.playerCards {
		if i == oldIdx {
			continue
		}
		newHand = append(newHand, card)
	}
	newHand = append(newHand, window.playerCards[oldIdx])
	return newHand
}

func (window *Window) insertBeforeIdx(oldIdx, newIdx int) []CardModel {
	newHand := make([]CardModel, 0)
	for i, card := range window.playerCards {
		if i == oldIdx {
			continue
		} else if i == newIdx {
			newHand = append(newHand, window.playerCards[oldIdx])
		}
		newHand = append(newHand, card)
	}
	return newHand
}

func (window *Window) cardModelsToCards() []*dm.Card {
	retVal := make([]*dm.Card, 0)
	for _, card := range window.playerCards {
		retVal = append(retVal, card.srcCard)
	}
	return retVal
}

func (window *Window) sendRearrangedHand() {
	cards := window.cardModelsToCards()
	window.updatePlayerHand(cards)
	actionMsg := cm.NewActionRearrangeCardsMessage(window.clientId, cards)
	window.sendActionCallback(actionMsg)
}
