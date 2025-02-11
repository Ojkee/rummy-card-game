package deck_manager

import "math/rand"

type CardQueue struct {
	Cards []*Card `json:"cards"`
}

func NewCardQueue() *CardQueue {
	return &CardQueue{
		Cards: make([]*Card, 0),
	}
}

func (cq *CardQueue) Pop() *Card {
	if len(cq.Cards) == 0 {
		return nil
	}
	var card *Card
	card, cq.Cards = cq.Cards[0], cq.Cards[1:]
	return card
}

func (cq *CardQueue) PopBack() *Card {
	if len(cq.Cards) == 0 {
		return nil
	}
	var card *Card
	card, cq.Cards = cq.Cards[len(cq.Cards)-1], cq.Cards[:len(cq.Cards)-1]
	return card
}

func (cq *CardQueue) SeekBack() *Card {
	if len(cq.Cards) == 0 {
		return nil
	}
	return cq.Cards[len(cq.Cards)-1]
}

func (cq *CardQueue) ShuffleExtend(newCards []*Card) {
	numAllCards := len(newCards)
	for i := range numAllCards {
		rIdx := rand.Int31n(int32(numAllCards))
		newCards[i], newCards[rIdx] = newCards[rIdx], newCards[i]
	}
	cq.Extend(newCards)
}

func (cq *CardQueue) Extend(newCards []*Card) {
	cq.Cards = append(cq.Cards, newCards...)
}

func (cq *CardQueue) Push(newCard *Card) {
	cq.Cards = append(cq.Cards, newCard)
}

func (cq *CardQueue) Empty() bool {
	return cq.Left() == 0
}

func (cq *CardQueue) Left() int {
	return len(cq.Cards)
}
