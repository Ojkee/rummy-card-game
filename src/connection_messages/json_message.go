package connection_messages

import (
	"strconv"
)

type MESSAGE_TYPE int

const (
	ID_INFO MESSAGE_TYPE = iota
	STATE_VIEW
	PLAYER_READY
	PLAYER_ACTION
	GAME_STATE_INFO
	GAME_WINDOW_TEXT
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
	Json() ([]byte, error)
}

type DefaultMessage struct {
	MessageType MESSAGE_TYPE `json:"message_type"`
}

type ClientMessage struct {
	DefaultMessage
	ClientId int `json:"client_id"`
}

func (dm *DefaultMessage) GetMessageType() MESSAGE_TYPE {
	return dm.MessageType
}
