package connection_messages

type AcceptMessage struct {
	DefaultMessage
	HasAccepted bool `json:"has_accepted"`
}

func NewAcceptMessage(status bool) *AcceptMessage {
	return &AcceptMessage{
		DefaultMessage: DefaultMessage{
			MessageType: PLAYER_ACCEPT,
		},
		HasAccepted: status,
	}
}
