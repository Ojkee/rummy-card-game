package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type SequenceModel struct {
	cardModels   []CardModel
	firstCardPos rl.Vector2
}

func NewSequenceModel(cards []*dm.Card, firstCardPos rl.Vector2) *SequenceModel {
	return &SequenceModel{
		cardModels:   CardsToModels(cards, firstCardPos, true),
		firstCardPos: firstCardPos,
	}
}

func CardsToModels(cards []*dm.Card, firstCardPos rl.Vector2, isSmall bool) []CardModel {
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
	var width float32 = 0
	var height float32 = 0
	return rl.NewVector2(width, height)
}

func (sm *SequenceModel) Draw() {
	for _, cardModel := range sm.cardModels {
		cardModel.Draw()
	}
}
