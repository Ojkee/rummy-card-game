package window

import (
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

func (window *Window) handleAvailableSpots(mousePos *rl.Vector2) {
	for _, availableSpot := range window.availableSpots {
		if availableSpot.InRect(mousePos) {
			selectedCard := window.getCardIfOneSelected()
			if selectedCard == nil {
				return
			}
			newUpdateSequenceMsg := cm.NewActionUpdateTableSequenceMessage(
				window.clientId,
				availableSpot.GetSequence().Id,
				window.getCardIdxByType(availableSpot.GetSpotType(), availableSpot.GetSequence()),
				selectedCard.srcCard,
			)
			window.sendActionCallback(newUpdateSequenceMsg)
			return
		}
	}
}

func (window *Window) getCardIdxByType(spotType gm.AVAILABLE_SPOT_TYPE, seq gm.Sequence) int {
	switch spotType {
	case gm.ADD_BEGIN:
		return -1
	case gm.ADD_END:
		return len(seq.TableCards) + 1
	}
	return 0
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
	for i, card := range window.playerCards {
		if card.isSelected && card.sequenceId == -1 {
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
	for i, cardModel := range window.playerCards {
		if cardModel.sequenceId == seqId {
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

func (window *Window) collectLockedSequencesCards() []*cm.SequenceLocked {
	retVal := make([]*cm.SequenceLocked, 0)
	for seqId, isUsed := range window.lockedSequencesIds {
		if isUsed {
			cardsInSameSeq := make([]*dm.Card, 0)
			for _, card := range window.playerCards {
				if card.sequenceId == seqId {
					cardsInSameSeq = append(cardsInSameSeq, card.srcCard)
				}
			}
			if len(cardsInSameSeq) > 0 {
				seqLocked := cm.NewSequenceLocked(seqId, cardsInSameSeq)
				retVal = append(retVal, seqLocked)
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
	window.availableSpots = make([]gm.AvailableSpot, 0)
	for _, sequenceModel := range window.tableSequences {
		ids := gm.FitSequenceIds(cardModel.srcCard, sequenceModel.sequence)
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
			availableSpot := gm.NewAvailableSpot(
				rect,
				gm.ADD_BEGIN,
				COLOR_HIGHLIGHT_SPOT,
				*sequenceModel.sequence,
			)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		} else if idx >= len(sequenceModel.cardModels) {
			rect := newRect(sequenceModel.firstCardPos.X + sequenceModel.GetSize().X)
			availableSpot := gm.NewAvailableSpot(rect, gm.ADD_END, COLOR_HIGHLIGHT_SPOT, *sequenceModel.sequence)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		} else {
			rect := newRect(sequenceModel.firstCardPos.X + float32(idx)*SEQUENCE_CARD_WIDTH)
			availableSpot := gm.NewAvailableSpot(rect, gm.REPLACE_JOKER, COLOR_HIGHLIGHT_SPOT, *sequenceModel.sequence)
			window.availableSpots = append(window.availableSpots, *availableSpot)
		}
	}
}

func (window *Window) resetAvailableSpots() {
	window.availableSpots = nil
	window.availableSpots = make([]gm.AvailableSpot, 0)
}
