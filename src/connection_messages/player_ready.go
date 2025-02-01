package connection_messages

type ReadyMessage struct {
	DefaultMessage
	IsReady bool `json:"is_ready"`
}

func NewReadyMessage(status bool) *ReadyMessage {
	return &ReadyMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PLAYER_READY,
		},
		IsReady: status,
	}
}
