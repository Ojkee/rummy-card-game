package connection_messages

import (
	"encoding/json"
)

type WrongCardsHighlight struct {
	DefaultMessage
	SeqLocked []*SequenceLocked `json:"sequence_locked"`
}

func NewWrongCardsHighlight(seqLocked []*SequenceLocked) *WrongCardsHighlight {
	return &WrongCardsHighlight{
		DefaultMessage: DefaultMessage{
			MessageType: WRONG_CARDS_HIGHLIGHT,
		},
		SeqLocked: seqLocked,
	}
}

func (wch *WrongCardsHighlight) Json() ([]byte, error) {
	return json.Marshal(wch)
}
