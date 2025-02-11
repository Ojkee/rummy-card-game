package connection_messages

import (
	"encoding/json"
)

type ACTION_TYPE int

const (
	DRAW_CARD ACTION_TYPE = iota
)

type ActionDrawMessage struct {
	DefaultMessage
	ActionType ACTION_TYPE
}

func NewActionDrawMessage() *ActionDrawMessage {
	return &ActionDrawMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PLAYER_ACTION,
		},
		ActionType: DRAW_CARD,
	}
}

func (adm *ActionDrawMessage) Json() ([]byte, error) {
	return json.Marshal(adm)
}
