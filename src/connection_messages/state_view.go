package connection_messages

import (
	"encoding/json"

	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
	"rummy-card-game/src/game_logic/player"
)

type StateView struct {
	DefaultMessage
	TurnPlayerId      int            `json:"turn_player_id"`
	DrawPile          *dm.CardQueue  `json:"draw_pile"`
	DiscardPile       *dm.CardQueue  `json:"discard_pile"`
	PlayerEntity      *player.Player `json:"player_entity"`
	OpponentsNumCards []int          `json:"opponents_num_cards"`
	TableSequences    []gm.Sequence  `json:"table_sequences"`
}

func NewStateView(
	turnPlayerId int,
	drawPile, discardPile *dm.CardQueue,
	playerEntity *player.Player,
	opponentsNumCards []int,
	tableSequences []gm.Sequence,
) *StateView {
	return &StateView{
		DefaultMessage: DefaultMessage{
			MessageType: STATE_VIEW,
		},
		TurnPlayerId:      turnPlayerId,
		DrawPile:          drawPile,
		DiscardPile:       discardPile,
		PlayerEntity:      playerEntity,
		OpponentsNumCards: opponentsNumCards,
		TableSequences:    tableSequences,
	}
}

func (sv *StateView) Json() ([]byte, error) {
	return json.Marshal(sv)
}
