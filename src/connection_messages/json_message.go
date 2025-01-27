package connection_messages

import "strconv"

type MESSAGE_TYPE int

const (
	ID_INFO MESSAGE_TYPE = iota
	STATE_VIEW
	PLAYER_ACCEPT
	PLAYER_ACTION
	UNKNOWN
)

func MessageTypeOfString(typeStr string) MESSAGE_TYPE {
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		return UNKNOWN
	}
	if typeInt < int(ID_INFO) || typeInt > int(UNKNOWN) {
		return UNKNOWN
	}
	return MESSAGE_TYPE(typeInt)
}

type JsonMessage interface {
	GetMessageType() MESSAGE_TYPE
}

type DefaultMessage struct {
	MessageType MESSAGE_TYPE `json:"message_type"`
}

func (dm *DefaultMessage) GetMessageType() MESSAGE_TYPE {
	return dm.MessageType
}
