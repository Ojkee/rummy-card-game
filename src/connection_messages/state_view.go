package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
	"rummy-card-game/src/game_logic/player"
)

type StateView struct {
	DefaultMessage
	DrawPile          *dm.CardQueue  `json:"draw_pile"`
	DiscardPile       *dm.CardQueue  `json:"discard_pile"`
	PlayerEntity      *player.Player `json:"player_entity"`
	OpponentsNumCards []int          `json:"opponents_num_cards"`
}

func NewStateView(
	drawPile, discardPile *dm.CardQueue,
	playerEntity *player.Player,
	opponentsNumCards []int,
) *StateView {
	return &StateView{
		DefaultMessage: DefaultMessage{
			MessageType: STATE_VIEW,
		},
		DrawPile:          drawPile,
		DiscardPile:       discardPile,
		PlayerEntity:      playerEntity,
		OpponentsNumCards: opponentsNumCards,
	}
}

func (sv *StateView) Json() ([]byte, error) {
	return json.Marshal(sv)
}
