package connection_messages

import "encoding/json"

type DEBUG_MESSAGE_TYPE int

const (
	DEBUG_RESET DEBUG_MESSAGE_TYPE = iota
	DEBUG_UNSUPPORTED
)

type DebugMessage interface {
	JsonMessage
	GetDebugMessageType() DEBUG_MESSAGE_TYPE
}

type ResetGameMessage struct {
	DefaultMessage
	DebugType DEBUG_MESSAGE_TYPE `json:"debug_type"`
}

func NewResetGameMessage() *ResetGameMessage {
	return &ResetGameMessage{
		DefaultMessage: DefaultMessage{
			MessageType: DEBUG_MESSAGE,
		},
		DebugType: DEBUG_RESET,
	}
}

func (rgm *ResetGameMessage) Json() ([]byte, error) {
	return json.Marshal(rgm)
}

func (rgm *ResetGameMessage) GetDebugMessageType() DEBUG_MESSAGE_TYPE {
	return rgm.DebugType
}
