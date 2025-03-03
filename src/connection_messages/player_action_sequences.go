package connection_messages

import (
	"encoding/json"
)

type ActionMeldMessage struct {
	ClientMessage
	ActionType ACTION_TYPE       `json:"action_type"`
	Sequences  []*SequenceLocked `json:"sequences"`
}

func NewActionMeldMessage(clientId int, sequences []*SequenceLocked) *ActionMeldMessage {
	return &ActionMeldMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: INITIAL_MELD,
		Sequences:  sequences,
	}
}

func (amm *ActionMeldMessage) Json() ([]byte, error) {
	return json.Marshal(amm)
}

func (amm *ActionMeldMessage) GetActionType() ACTION_TYPE {
	return amm.ActionType
}
