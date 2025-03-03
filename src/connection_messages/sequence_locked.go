package connection_messages

import dm "rummy-card-game/src/game_logic/deck_manager"

type SequenceLocked struct {
	SequenceId int        `json:"sequence_id"`
	Cards      []*dm.Card `json:"cards"`
}

func NewSequenceLocked(seqId int, cards []*dm.Card) *SequenceLocked {
	return &SequenceLocked{
		SequenceId: seqId,
		Cards:      cards,
	}
}
