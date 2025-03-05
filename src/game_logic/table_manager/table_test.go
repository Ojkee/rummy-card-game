package table_manager

import (
	"fmt"
	"testing"

	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

type tuple struct {
	a, b *dm.Card
}

func zip(a, b []*dm.Card) []tuple {
	r := make([]tuple, len(a), len(a))
	for i, e := range a {
		r[i] = tuple{e, b[i]}
	}
	return r
}

func seqsAreTheSame(resultSeq, expectedSeq []*dm.Card, t *testing.T) bool {
	if len(resultSeq) != len(expectedSeq) {
		t.Errorf("seqs lenghts aren't the same")
		return false
	}
	for _, row := range zip(resultSeq, expectedSeq) {
		if *row.a != *row.b {
			return false
		}
	}
	return true
}

func jokImitationsAreTheSame(resultImit, expectedImit []gm.JokerImitation, t *testing.T) bool {
	if len(resultImit) != len(expectedImit) {
		t.Errorf("jok imitations lenghts aren't the same")
		return false
	}
	for i := range resultImit {
		if *resultImit[i].Card != *expectedImit[i].Card ||
			resultImit[i].Idx != expectedImit[i].Idx {
			return false
		}
	}
	return true
}

func TestSortAscendingSeq(t *testing.T) {
	table := Table{}
	tests := []struct {
		inputSeq         []*dm.Card
		expectedSeq      []*dm.Card
		expectedJokImits []gm.JokerImitation
	}{
		{ // 0
			inputSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.THREE),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.FOUR),
				dm.NewCard(dm.CLUBS, dm.FIVE),
			},
			expectedJokImits: []gm.JokerImitation{},
		},
		{ // 1
			inputSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.THREE),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.CLUBS, dm.FOUR)),
			},
		},
		{ // 2
			inputSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.CLUBS, dm.FOUR)),
				*gm.NewJokerImitation(4, dm.NewCard(dm.CLUBS, dm.SEVEN)),
			},
		},
		{ // 3
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
				dm.NewCard(dm.CLUBS, dm.THREE),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.CLUBS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.CLUBS, dm.FIVE),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.CLUBS, dm.FOUR)),
			},
		},
		{ // 4
			inputSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.KING),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(2, dm.NewCard(dm.DIAMONDS, dm.ACE)),
			},
		},
		{ // 5
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.JACK)),
			},
		},
		{ // 6
			inputSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.ACE),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.JACK)),
			},
		},
		{ // 7
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.TEN)),
			},
		},
		{ // 8
			inputSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(3, dm.NewCard(dm.DIAMONDS, dm.ACE)),
			},
		},
		{ // 9
			inputSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.ACE),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.TEN)),
			},
		},
		{ // 10
			inputSeq: []*dm.Card{
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.TEN)),
				*gm.NewJokerImitation(4, dm.NewCard(dm.DIAMONDS, dm.ACE)),
			},
		},
		{ // 11
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.NINE)),
				*gm.NewJokerImitation(1, dm.NewCard(dm.DIAMONDS, dm.TEN)),
			},
		},
		{ // 12
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.DIAMONDS, dm.JACK),
				dm.NewCard(dm.DIAMONDS, dm.QUEEN),
				dm.NewCard(dm.DIAMONDS, dm.KING),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.DIAMONDS, dm.TEN)),
				*gm.NewJokerImitation(4, dm.NewCard(dm.DIAMONDS, dm.ACE)),
			},
		},
		{ // 13
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.THREE),
				dm.NewCard(dm.HEARTS, dm.TWO),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.TWO),
				dm.NewCard(dm.HEARTS, dm.THREE),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.HEARTS, dm.ACE)),
			},
		},
		{ // 14
			inputSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.TWO),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.TWO),
				dm.NewCard(dm.HEARTS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(2, dm.NewCard(dm.HEARTS, dm.FOUR)),
			},
		},
		{ // 15
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.THREE),
				dm.NewCard(dm.HEARTS, dm.TWO),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.TWO),
				dm.NewCard(dm.HEARTS, dm.THREE),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.HEARTS, dm.ACE)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.FOUR)),
			},
		},
		{ // 16
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SEVEN),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.SEVEN),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.HEARTS, dm.FIVE)),
			},
		},
		{ // 17
			inputSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.SEVEN),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.SEVEN),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.NINE)),
			},
		},
		{ // 18
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.SEVEN)),
			},
		},
		{ // 19
			inputSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.SEVEN)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.NINE)),
				*gm.NewJokerImitation(4, dm.NewCard(dm.HEARTS, dm.TEN)),
			},
		},
		{ // 20
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.SIX),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.HEARTS, dm.FOUR)),
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.FIVE)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.SEVEN)),
			},
		},
		{ // 21
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.FOUR),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.FIVE)),
				*gm.NewJokerImitation(2, dm.NewCard(dm.HEARTS, dm.SIX)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.SEVEN)),
			},
		},
		{ // 22
			inputSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.FIVE)),
				*gm.NewJokerImitation(2, dm.NewCard(dm.HEARTS, dm.SIX)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.SEVEN)),
			},
		},
		{ // 23
			inputSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.ANY, dm.JOKER),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(1, dm.NewCard(dm.HEARTS, dm.FIVE)),
				*gm.NewJokerImitation(2, dm.NewCard(dm.HEARTS, dm.SIX)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.SEVEN)),
				*gm.NewJokerImitation(5, dm.NewCard(dm.HEARTS, dm.NINE)),
			},
		},
		{ // 24
			inputSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
				dm.NewCard(dm.HEARTS, dm.FOUR),
			},
			expectedSeq: []*dm.Card{
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.FOUR),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.ANY, dm.JOKER),
				dm.NewCard(dm.HEARTS, dm.EIGHT),
			},
			expectedJokImits: []gm.JokerImitation{
				*gm.NewJokerImitation(0, dm.NewCard(dm.HEARTS, dm.THREE)),
				*gm.NewJokerImitation(2, dm.NewCard(dm.HEARTS, dm.FIVE)),
				*gm.NewJokerImitation(3, dm.NewCard(dm.HEARTS, dm.SIX)),
				*gm.NewJokerImitation(4, dm.NewCard(dm.HEARTS, dm.SEVEN)),
			},
		},
	}

	passed := make([]int, 0)
	testRange := tests[:]
	for i, test := range testRange {
		errored := false
		result, resultImitations := table.sortAscendingSequence(test.inputSeq)
		if !seqsAreTheSame(result, test.expectedSeq, t) {
			errored = true
			t.Errorf(
				"TEST %d:\ninput: %v\n\tresult: %v\n\texpected: %v",
				i,
				test.inputSeq,
				result,
				test.expectedSeq,
			)
		}
		if !jokImitationsAreTheSame(resultImitations, test.expectedJokImits, t) {
			errored = true
			t.Errorf(
				"TEST %d:\ninput: %v\n\tresultImit: %v\n\texpectedImit: %v",
				i,
				test.inputSeq,
				resultImitations,
				test.expectedJokImits,
			)
		}

		if !errored {
			passed = append(passed, i)
		}
	}
	if len(passed) != len(tests) {
		fmt.Printf("Passed: %d/%d\n", len(passed), len(testRange))
	}
}
