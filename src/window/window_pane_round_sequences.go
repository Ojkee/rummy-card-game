package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
)

func (window *Window) handleLockSequence(mousePos *rl.Vector2) {
	if selectedCards := window.getSelectedUnlockedCards(); len(selectedCards) >= 3 {
		if window.lockSetButton.isClicked(mousePos) {
			window.lockSelectedSequence()
		}
	}
}

func (window *Window) getSelectedUnlockedCards() []*CardModel {
	selectedCards := make([]*CardModel, 0)
	for _, card := range window.playerCards {
		if card.isSelected && card.sequenceId == -1 {
			selectedCards = append(selectedCards, &card)
		}
	}
	return selectedCards
}

func (window *Window) getNextLockId() int {
	for key, val := range window.lockedSequencesIds {
		if !val {
			window.lockedSequencesIds[key] = true
			return key
		}
	}
	return -1
}

func (window *Window) lockSelectedSequence() {
	nextLockId := window.getNextLockId()
	for i := 0; i < len(window.playerCards); i++ {
		if window.playerCards[i].isSelected && window.playerCards[i].sequenceId == -1 {
			window.playerCards[i].SetSequenceId(nextLockId)
		}
	}
}

func (window *Window) unlockAllById(seqId int) {
	for i := 0; i < len(window.playerCards); i++ {
		if window.playerCards[i].sequenceId == seqId {
			window.playerCards[i].Reset()
		}
	}
	window.lockedSequencesIds[seqId] = false
}

func (window *Window) numLockedSequences() int {
	lockCounter := 0
	for _, val := range window.lockedSequencesIds {
		if val {
			lockCounter++
		}
	}
	return lockCounter
}

func (window *Window) collectLockedSequencesCards() [][]*dm.Card {
	retVal := make([][]*dm.Card, 0)
	for seqId, isUsed := range window.lockedSequencesIds {
		if isUsed {
			cardsInSameSeq := make([]*dm.Card, 0)
			for _, card := range window.playerCards {
				if card.sequenceId == seqId {
					cardsInSameSeq = append(cardsInSameSeq, card.srcCard)
				}
			}
			retVal = append(retVal, cardsInSameSeq)
		}
	}
	return retVal
}

func (window *Window) handleInitialMeldButton(mousePos *rl.Vector2) {
	if window.numLockedSequences() > 0 &&
		window.initialMeldButton.isClicked(mousePos) {
		lockedSequences := window.collectLockedSequencesCards()
		lockedSequencesMessage := cm.NewActionInitialMeldMessage(window.clientId, lockedSequences)
		window.sendActionCallback(lockedSequencesMessage)
	}
}
