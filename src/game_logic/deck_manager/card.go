package deck_manager

import "fmt"

type (
	Suit int
	Rank int
)

const (
	SPADES Suit = iota
	HEARTS
	DIAMONDS
	CLUBS
	ANY // If joker
)

const (
	TWO Rank = iota
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	JACK
	QUEEN
	KING
	ACE
	JOKER
)

var (
	suits = []Suit{SPADES, HEARTS, DIAMONDS, CLUBS, ANY}
	ranks = []Rank{
		TWO,
		THREE,
		FOUR,
		FIVE,
		SIX,
		SEVEN,
		EIGHT,
		NINE,
		TEN,
		JACK,
		QUEEN,
		KING,
		ACE,
		JOKER,
	}
)

func SuitOfInt(n int) (Suit, error) {
	if n < 0 || n >= len(suits) {
		return 0, fmt.Errorf("invalid Suit: %d", n)
	}
	return suits[n], nil
}

func RankOfInt(n int) (Rank, error) {
	if n < 0 || n >= len(ranks) {
		return 0, fmt.Errorf("invalid Rank: %d", n)
	}
	return ranks[n], nil
}

func (s Suit) String() string {
	suits := []string{"SPADES", "HEARTS", "DIAMONDS", "CLUBS", "ANY"}
	return suits[s]
}

func (r Rank) String() string {
	ranks := []string{
		"2", "3", "4",
		"5", "6", "7",
		"8", "9", "10",
		"J", "Q", "K",
		"A", "JOKER",
	}
	return ranks[r]
}

type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rand"`
}

func NewCard(suit Suit, rank Rank) *Card {
	return &Card{
		Suit: suit,
		Rank: rank,
	}
}

func (card *Card) String() string {
	return fmt.Sprintf(
		"{%s %s}",
		card.Suit.String(),
		card.Rank.String(),
	)
}
