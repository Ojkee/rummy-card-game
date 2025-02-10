package connection_messages

type ReadyMessage struct {
	ClientMessage
	IsReady bool `json:"is_ready"`
}

func NewReadyMessage(status bool, ClientId int) *ReadyMessage {
	return &ReadyMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_READY,
			},
			ClientId: ClientId,
		},
		IsReady: status,
	}
}
