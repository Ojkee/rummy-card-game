package table_manager

import (
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
		t.Errorf("lenghts aren't the same")
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
		t.Errorf("lenghts aren't the same")
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
	}

	for i, test := range tests {
		result, resultImitations := table.sortAscendingSequence(test.inputSeq)
		if !seqsAreTheSame(result, test.expectedSeq, t) {
			t.Errorf(
				"TEST %d:\n\tresult: %v\n\texpected: %v",
				i,
				result,
				test.expectedSeq,
			)
		} else if !jokImitationsAreTheSame(resultImitations, test.expectedJokImits, t) {
			t.Errorf(
				"TEST %d:\n\tresultImit: %v\n\texpectedImit: %v",
				i,
				resultImitations,
				test.expectedJokImits,
			)
		}
	}
}
