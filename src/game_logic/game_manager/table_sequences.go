package game_manager

import dm "rummy-card-game/src/game_logic/deck_manager"

type SEQUENCE_TYPE int

const (
	SEQUENCE_SAME_RANK SEQUENCE_TYPE = iota
	SEQUENCE_ASCENDING
	SEQUENCE_PURE
)

type Sequence struct {
	TableCards []*dm.Card    `json:"table_cards"`
	Type       SEQUENCE_TYPE `json:"type"`
}

func NewSequence(cards []*dm.Card, sequenceType SEQUENCE_TYPE) *Sequence {
	return &Sequence{
		TableCards: cards,
		Type:       sequenceType,
	}
}
