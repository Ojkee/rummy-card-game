package window

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

func (window *Window) handleLockSequence(mousePos *rl.Vector2) {
	if selectedCards := window.getSelectedUnlockedCards(); len(selectedCards) >= 3 {
		if window.lockSetButton.InRect(mousePos) {
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

func (window *Window) resetLockedSequencesIds() {
	for key := range window.lockedSequencesIds {
		window.lockedSequencesIds[key] = false
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
	for _, isUsed := range window.lockedSequencesIds {
		if isUsed {
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
			if len(cardsInSameSeq) > 0 {
				retVal = append(retVal, cardsInSameSeq)
			}
		}
	}
	return retVal
}

func (window *Window) handleInitialMeldButton(mousePos *rl.Vector2) {
	if window.numLockedSequences() > 0 &&
		window.initialMeldButton.InRect(mousePos) {
		lockedSequences := window.collectLockedSequencesCards()
		lockedSequencesMessage := cm.NewActionMeldMessage(window.clientId, lockedSequences)
		window.sendActionCallback(lockedSequencesMessage)
		window.resetLockedSequencesIds()
	}
}

func (window *Window) initAvailableSpots(cardModel *CardModel) {
	window.availableSpots = nil
	window.availableSpots = make([]AvailableSpot, 0)
	for _, sequenceModel := range window.tableSequences {
		ids := gm.FitSequenceIds(cardModel.srcCard, sequenceModel.sequence)
		log.Println(ids)
		window.addNewAvailableSpots(ids, sequenceModel)
	}
}

func (window *Window) addNewAvailableSpots(ids []int, sequenceModel SequenceModel) {
	newRect := func(x float32) rl.Rectangle {
		return rl.NewRectangle(
			x,
			sequenceModel.firstCardPos.Y,
			SEQUENCE_CARD_WIDTH,
			SEQUENCE_CARD_HEIGHT,
		)
	}
	for _, idx := range ids {
		if idx < 0 {
			rect := newRect(sequenceModel.firstCardPos.X - SEQUENCE_CARD_WIDTH)
			availableSpot := NewAvailableSpot(rect, ADD_BEGIN)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		} else if idx >= len(sequenceModel.cardModels) {
			rect := newRect(sequenceModel.firstCardPos.X + sequenceModel.GetSize().X)
			availableSpot := NewAvailableSpot(rect, ADD_END)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		} else {
			rect := newRect(sequenceModel.firstCardPos.X + float32(idx)*SEQUENCE_CARD_WIDTH)
			availableSpot := NewAvailableSpot(rect, REPLACE_JOKER)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		}
	}
}

func (window *Window) resetAvailableSpots() {
	window.availableSpots = nil
	window.availableSpots = make([]AvailableSpot, 0)
}
