package connection_messages

import "encoding/json"

type IdInfo struct {
	DefaultMessage
	Id int `json:"id"`
}

func NewIdInfo(id int) *IdInfo {
	return &IdInfo{
		DefaultMessage: DefaultMessage{
			MessageType: ID_INFO,
		},
		Id: id,
	}
}

func (ii *IdInfo) Json() ([]byte, error) {
	return json.Marshal(ii)
}
