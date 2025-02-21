package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type ActionInitialMeldMessage struct {
	ClientMessage
	ActionType ACTION_TYPE  `json:"action_type"`
	Sequences  [][]*dm.Card `json:"sequences"`
}

func NewActionInitialMeldMessage(clientId int, sequences [][]*dm.Card) *ActionInitialMeldMessage {
	return &ActionInitialMeldMessage{
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

func (aimm *ActionInitialMeldMessage) Json() ([]byte, error) {
	return json.Marshal(aimm)
}

func (aimm *ActionInitialMeldMessage) GetActionType() ACTION_TYPE {
	return aimm.ActionType
}
