package network_server

import "github.com/gorilla/websocket"

type ConnectedClient struct {
	isReady  bool
	nickname string
	conn     *websocket.Conn

	// in round
	drawnCard bool
	hasMelded bool
}

func NewConnectedClient(conn *websocket.Conn) *ConnectedClient {
	return &ConnectedClient{
		isReady: false,
		conn:    conn,

		drawnCard: false,
	}
}

func (cc *ConnectedClient) AfterRoundReset() {
	cc.drawnCard = false
}

func (cc *ConnectedClient) Reset() {
	cc.drawnCard = false
	cc.hasMelded = false
}
