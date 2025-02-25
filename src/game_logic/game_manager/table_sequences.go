package game_manager

import dm "rummy-card-game/src/game_logic/deck_manager"

type SEQUENCE_TYPE int

const (
	SEQUENCE_SAME_RANK SEQUENCE_TYPE = iota
	SEQUENCE_ASCENDING
	SEQUENCE_PURE
)

type Sequence struct {
	TableCards      []*dm.Card         `json:"table_cards"`
	Type            SEQUENCE_TYPE      `json:"type"`
	JokerImitations map[string]dm.Card `json:"joker_imitations"`
}

func NewSequence(cards []*dm.Card, sequenceType SEQUENCE_TYPE) *Sequence {
	return &Sequence{
		TableCards:      cards,
		Type:            sequenceType,
		JokerImitations: make(map[string]dm.Card),
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

func (s *Sequence) SetJokerImitations(jokerImitations map[string]dm.Card) {
	s.JokerImitations = jokerImitations
}
