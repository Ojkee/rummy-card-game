package connection_messages

import "encoding/json"

func DecodeMessageType(msg []byte) (MESSAGE_TYPE, error) {
	var messageDecoded map[string]json.RawMessage
	var err error
	err = json.Unmarshal(msg, &messageDecoded)
	if err != nil {
		return UNKNOWN, err
	}
	var messageType MESSAGE_TYPE
	err = json.Unmarshal(messageDecoded["message_type"], &messageType)
	if err != nil {
		return UNKNOWN, err
	}
	return messageType, nil
}

func DecodeMessageClientId(msg []byte) (int, error) {
	var messageDecoded map[string]json.RawMessage
	var err error
	err = json.Unmarshal(msg, &messageDecoded)
	if err != nil {
		return -1, err
	}
	var clientId int
	err = json.Unmarshal(messageDecoded["client_id"], &clientId)
	if err != nil {
		return -1, err
	}
	return clientId, nil
}
