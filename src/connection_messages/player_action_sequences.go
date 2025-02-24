package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type ActionMeldMessage struct {
	ClientMessage
	ActionType ACTION_TYPE  `json:"action_type"`
	Sequences  [][]*dm.Card `json:"sequences"`
}

func NewActionMeldMessage(clientId int, sequences [][]*dm.Card) *ActionMeldMessage {
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
