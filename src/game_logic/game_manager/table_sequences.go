package game_manager

import dm "rummy-card-game/src/game_logic/deck_manager"

type SEQUENCE_TYPE int

const (
	SEQUENCE_SAME_RANK SEQUENCE_TYPE = iota
	SEQUENCE_ASCENDING
	SEQUENCE_PURE
)

type JokerImitation struct {
	Idx  int      `json:"idx"`
	Card *dm.Card `json:"imit_card"`
}

func NewJokerImitation(idx int, card *dm.Card) *JokerImitation {
	return &JokerImitation{
		Idx:  idx,
		Card: card,
	}
}

type Sequence struct {
	Id              int              `json:"id"`
	TableCards      []*dm.Card       `json:"table_cards"`
	Type            SEQUENCE_TYPE    `json:"type"`
	JokerImitations []JokerImitation `json:"joker_imitations"`
}

func NewSequence(
	id int,
	cards []*dm.Card,
	sequenceType SEQUENCE_TYPE,
	jokerImitations []JokerImitation,
) *Sequence {
	return &Sequence{
		Id:              id,
		TableCards:      cards,
		Type:            sequenceType,
		JokerImitations: jokerImitations,
	}
}

func (s *Sequence) GetSuitIfAscending() dm.Suit {
	for _, card := range s.TableCards {
		if card.Suit != dm.ANY {
			return card.Suit
		}
	}
	return dm.ANY
}

func (s *Sequence) GetId() int {
	return s.Id
}
