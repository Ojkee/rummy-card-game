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
			t.Errorf("Error in test %d: got %v; want: %v", i, result, test.expected)
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
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.SPADES, dm.FIVE),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.SPADES, dm.FIVE),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.JACK),
				dm.NewCard(dm.CLUBS, dm.QUEEN),
				dm.NewCard(dm.CLUBS, dm.KING),
				dm.NewCard(dm.CLUBS, dm.ACE),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.CLUBS, dm.SEVEN),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.NINE),
				dm.NewCard(dm.CLUBS, dm.TEN),
				dm.NewCard(dm.CLUBS, dm.JACK),
				dm.NewCard(dm.CLUBS, dm.QUEEN),
				dm.NewCard(dm.CLUBS, dm.KING),
				dm.NewCard(dm.CLUBS, dm.ACE),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.CLUBS, dm.SEVEN),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.NINE),
				dm.NewCard(dm.CLUBS, dm.TEN),
				dm.NewCard(dm.CLUBS, dm.JACK),
				dm.NewCard(dm.CLUBS, dm.QUEEN),
				dm.NewCard(dm.CLUBS, dm.KING),
				dm.NewCard(dm.CLUBS, dm.ACE),
				dm.NewCard(dm.CLUBS, dm.JOKER),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.CLUBS, dm.SEVEN),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.NINE),
				dm.NewCard(dm.CLUBS, dm.TEN),
				dm.NewCard(dm.CLUBS, dm.JACK),
				dm.NewCard(dm.CLUBS, dm.KING),
				dm.NewCard(dm.CLUBS, dm.ACE),
				dm.NewCard(dm.CLUBS, dm.JOKER),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.CLUBS, dm.SEVEN),
				dm.NewCard(dm.CLUBS, dm.EIGHT),
				dm.NewCard(dm.CLUBS, dm.NINE),
				dm.NewCard(dm.CLUBS, dm.TEN),
				dm.NewCard(dm.CLUBS, dm.JACK),
				dm.NewCard(dm.CLUBS, dm.KING),
				dm.NewCard(dm.CLUBS, dm.ACE),
			},
			false,
		},
		// {
		// 	[]*dm.Card{
		// 		dm.NewCard(dm.CLUBS, dm.ACE),
		// 		dm.NewCard(dm.CLUBS, dm.TWO),
		// 		dm.NewCard(dm.CLUBS, dm.THREE),
		// 	},
		// 	true,
		// },
	}

	for i, test := range tests {
		result := gm.IsAscendingSequence(test.sequence)
		if result != test.expected {
			t.Errorf("Error in test %d: got %v; want: %v", i, result, test.expected)
		}
	}
}

func TestIsPureSequence(t *testing.T) {
	tests := []struct {
		sequence []*dm.Card
		expected bool
	}{
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FIVE),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			true,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.DIAMONDS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			false,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.CLUBS, dm.TWO),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
			},
			true,
		},
	}

	for i, test := range tests {
		result := gm.IsPureSequence(test.sequence)
		if result != test.expected {
			t.Errorf("Error in test %d: got %v; want: %v", i, result, test.expected)
		}
	}
}

func TestPointsCounter(t *testing.T) {
	tests := []struct {
		sequence []*dm.Card
		expected int
	}{
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.TWO),
				dm.NewCard(dm.DIAMONDS, dm.THREE),
				dm.NewCard(dm.DIAMONDS, dm.FOUR),
			},
			9,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.ACE),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
			},
			33,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
			},
			31,
		},
		{
			[]*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			0,
		},
	}

	for i, test := range tests {
		result := gm.SequencePoints(test.sequence)
		if result != test.expected {
			t.Errorf("Error in test %d: got %v; want: %v", i, result, test.expected)
		}
	}
}
