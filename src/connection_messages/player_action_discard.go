package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type ActionDiscardMessage struct {
	ClientMessage
	ActionType ACTION_TYPE `json:"action_type"`
	Card       *dm.Card    `json:"card_id"`
}

func NewActionDiscardMessage(clientId int, card *dm.Card) *ActionDiscardMessage {
	return &ActionDiscardMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: DISCARD_CARD,
		Card:       card,
	}
}

func (adm *ActionDiscardMessage) Json() ([]byte, error) {
	return json.Marshal(adm)
}

func (adm *ActionDiscardMessage) GetActionType() ACTION_TYPE {
	return adm.ActionType
}
