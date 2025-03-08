package connection_messages

type ReadyMessage struct {
	ClientMessage
	IsReady  bool   `json:"is_ready"`
	Nickname string `json:"nickname"`
}

func NewReadyMessage(status bool, nickname string, ClientId int) *ReadyMessage {
	return &ReadyMessage{
		ClientMessage: ClientMessage{
			DefaultMessage: DefaultMessage{
				MessageType: PLAYER_READY,
			},
			ClientId: ClientId,
		},
		IsReady:  status,
		Nickname: nickname,
	}
}
