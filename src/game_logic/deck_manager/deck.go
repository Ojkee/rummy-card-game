package deck_manager

import "fmt"

type Deck struct {
	cards    []Card
	numCards int
}

func NewDeck() *Deck {
	_cards := make([]Card, 0)
	const numSuits = 4
	const numRanks = 13
	for i := range numSuits {
		for j := range numRanks {
			suit, _ := SuitOfInt(i)
			rank, _ := RankOfInt(j)
			_cards = append(_cards, *NewCard(suit, rank))
		}
	}
	_cards = append(_cards, *NewCard(ANY, JOKER))
	_cards = append(_cards, *NewCard(ANY, JOKER))
	return &Deck{
		cards:    _cards,
		numCards: len(_cards),
	}
}

func (deck *Deck) Print() {
	for i, card := range deck.cards {
		if i%4 == 0 {
			fmt.Println("")
		}
		fmt.Printf("%s ", card.String())
	}
}

func (deck *Deck) GetCards() *[]Card {
	return &deck.cards
}

func (deck *Deck) GetNumCards() int {
	return deck.numCards
}
