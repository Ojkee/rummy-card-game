package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

type SequenceModel struct {
	sequence     *gm.Sequence
	cardModels   []CardModel
	firstCardPos rl.Vector2
}

func NewSequenceModel(sequence *gm.Sequence, firstCardPos rl.Vector2) *SequenceModel {
	return &SequenceModel{
		sequence:     sequence,
		cardModels:   cardsToModels(sequence.TableCards, firstCardPos, true),
		firstCardPos: firstCardPos,
	}
}

func cardsToModels(cards []*dm.Card, firstCardPos rl.Vector2, isSmall bool) []CardModel {
	cardModels := make([]CardModel, 0)
	for i, card := range cards {
		x := firstCardPos.X + float32(i*int(SEQUENCE_CARD_WIDTH))
		y := firstCardPos.Y
		cardRect := rl.NewRectangle(
			x,
			y,
			SEQUENCE_CARD_WIDTH,
			SEQUENCE_CARD_HEIGHT,
		)
		cardModel := NewCardModel(card, cardRect, isSmall)
		cardModels = append(cardModels, *cardModel)
	}
	return cardModels
}

func (sm *SequenceModel) GetFirstCardPos() *rl.Vector2 {
	return &sm.firstCardPos
}

func (sm *SequenceModel) GetSize() rl.Vector2 {
	width := float32(len(sm.cardModels) * int(SEQUENCE_CARD_WIDTH))
	return rl.NewVector2(width, SEQUENCE_CARD_HEIGHT)
}

func (sm *SequenceModel) Draw() {
	for _, cardModel := range sm.cardModels {
		cardModel.Draw()
	}
}

func (sm *SequenceModel) GetSrcCards() []*dm.Card {
	srcCards := make([]*dm.Card, len(sm.cardModels))
	for i := range sm.cardModels {
		srcCards[i] = sm.cardModels[i].srcCard
	}
	return srcCards
}
