package connection_messages

import (
	"encoding/json"

	"rummy-card-game/src/game_logic/game_manager"
)

type GameStateInfo struct {
	DefaultMessage
	GameStateValue game_manager.GAME_STATE `json:"game_stete_value"`
}

func NewGameStateInfo(gameStateValue game_manager.GAME_STATE) *GameStateInfo {
	return &GameStateInfo{
		DefaultMessage: DefaultMessage{
			MessageType: GAME_STATE_INFO,
		},
		GameStateValue: gameStateValue,
	}
}

func (gmi *GameStateInfo) Json() ([]byte, error) {
	return json.Marshal(gmi)
}
