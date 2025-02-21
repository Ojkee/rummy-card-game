package connection_messages

import (
	"encoding/json"
)

type ActionDrawMessage struct {
	ClientMessage
	ActionType ACTION_TYPE `json:"action_type"`
}

func NewActionDrawMessage(clientId int) *ActionDrawMessage {
	return &ActionDrawMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: DRAW_CARD,
	}
}

func (adm *ActionDrawMessage) Json() ([]byte, error) {
	return json.Marshal(adm)
}

func (adm *ActionDrawMessage) GetActionType() ACTION_TYPE {
	return adm.ActionType
}
