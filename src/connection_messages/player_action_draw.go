package connection_messages

import (
	"encoding/json"
)

type DRAW_TYPE int

const (
	DRAW_FROM_PILE DRAW_TYPE = iota
	DRAW_FROM_DISCARD_PILE
)

type ActionDrawMessage struct {
	ClientMessage
	ActionType ACTION_TYPE `json:"action_type"`
	DrawSource DRAW_TYPE   `json:"draw_source"`
}

func NewActionDrawMessage(clientId int, drawSource DRAW_TYPE) *ActionDrawMessage {
	return &ActionDrawMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: DRAW_CARD,
		DrawSource: drawSource,
	}
}

func (adm *ActionDrawMessage) Json() ([]byte, error) {
	return json.Marshal(adm)
}

func (adm *ActionDrawMessage) GetActionType() ACTION_TYPE {
	return adm.ActionType
}
