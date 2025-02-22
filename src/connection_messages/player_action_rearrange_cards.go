package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type ActionRearrangeCardsMessage struct {
	ClientMessage
	ActionType ACTION_TYPE `json:"action_type"`
	Cards      []*dm.Card  `json:"sequences"`
}

func NewActionRearrangeCardsMessage(clientId int, cards []*dm.Card) *ActionRearrangeCardsMessage {
	return &ActionRearrangeCardsMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: REARRANGE_CARDS,
		Cards:      cards,
	}
}

func (arcm *ActionRearrangeCardsMessage) Json() ([]byte, error) {
	return json.Marshal(arcm)
}

func (arcm *ActionRearrangeCardsMessage) GetActionType() ACTION_TYPE {
	return arcm.ActionType
}
