package connection_messages

import "encoding/json"

type GameWindowText struct {
	DefaultMessage
	Value string `json:"value"`
}

func NewGameWindowText(textMsg string) *GameWindowText {
	return &GameWindowText{
		DefaultMessage: DefaultMessage{
			MessageType: GAME_WINDOW_TEXT,
		},
		Value: textMsg,
	}
}

func (gwt *GameWindowText) Json() ([]byte, error) {
	return json.Marshal(gwt)
}
