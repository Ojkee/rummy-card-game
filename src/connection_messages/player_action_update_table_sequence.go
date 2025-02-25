package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type ActionUpdateTableSequenceMessage struct {
	ClientMessage
	ActionType ACTION_TYPE `json:"action_type"`
	SequenceId int         `json:"sequence_id"`
	CardIdx    int         `json:"card_idx"`
	Card       *dm.Card
}

func NewActionUpdateTableSequenceMessage(
	clientId int,
	sequenceId, cardIdx int,
	card *dm.Card,
) *ActionUpdateTableSequenceMessage {
	return &ActionUpdateTableSequenceMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_ACTION,
			},
			ClientId: clientId,
		},
		ActionType: UPDATE_TABLE_SEQUNCE,
		SequenceId: sequenceId,
		CardIdx:    cardIdx,
		Card:       card,
	}
}

func (autsm *ActionUpdateTableSequenceMessage) Json() ([]byte, error) {
	return json.Marshal(autsm)
}

func (autsm *ActionUpdateTableSequenceMessage) GetActionType() ACTION_TYPE {
	return autsm.ActionType
}
