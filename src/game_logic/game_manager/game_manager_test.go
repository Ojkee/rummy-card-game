package game_manager_test

import (
	"testing"

	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

func TestIsSameRankSequence(t *testing.T) {
	tests := []struct {
		sequence []*dm.Card
		expected bool
	}{
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.EIGHT),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.EIGHT),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.SEVEN),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.SPADES, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			false,
		},
	}

	for i, test := range tests {
		result := gm.IsSameRankSequence(test.sequence)
		if result != test.expected {
			t.Errorf("Error in test: %d: got %v; want: %v", i, result, test.expected)
		}
	}
}

func TestIsSameAscendingSequence(t *testing.T) {
	tests := []struct {
		sequence []*dm.Card
		expected bool
	}{
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			true,
		},
	}

	for i, test := range tests {
		result := gm.IsAscendingSequence(test.sequence)
		if result != test.expected {
			t.Errorf("Error in test: %d: got %v; want: %v", i, result, test.expected)
		}
	}
}
