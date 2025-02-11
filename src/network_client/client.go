package network_client

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	"rummy-card-game/src/connection_messages"
	"rummy-card-game/src/window"
)

type Client struct {
	conn        *websocket.Conn
	id          int
	gameStarted bool
	gameWindow  window.Window
}

func NewClient() *Client {
	return &Client{
		conn:        nil,
		id:          -1,
		gameStarted: false,
		gameWindow:  *window.NewWindow(),
	}
}

func (client *Client) GetId() int {
	return client.id
}

func (client *Client) SetId(id int) {
	client.id = id
}

func (client *Client) Connect() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	client.conn = conn
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	defer conn.Close()
	log.Println("Connected!")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	client.gameWindow.SetOnReadyCallback(client.sendOnReady)

	go client.readFromServer(conn)
	client.gameWindow.MainLoop()

	select {
	case <-signalChannel:
		return
	case <-client.gameWindow.CloseListener():
		return
	}
}

func (client *Client) readFromServer(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Println("websocket closed: ", err)
			} else {
				log.Println("Unexpeced websocket error: ", err)
			}
			return
		}

		var messageType connection_messages.MESSAGE_TYPE
		messageType, err = connection_messages.DecodeMessageType(msg)
		if err != nil {
			log.Println("Error decoding message type:", err)
			continue
		}

		switch messageType {
		case connection_messages.ID_INFO:
			var idInfo connection_messages.IdInfo
			if err = json.Unmarshal(msg, &idInfo); err != nil {
				log.Println("Err parsing Id")
				continue
			}
			client.SetId(idInfo.Id)
		case connection_messages.STATE_VIEW:
			var stateView connection_messages.StateView
			if err = json.Unmarshal(msg, &stateView); err != nil {
				log.Println("Err parsing StateView")
				continue
			}
			client.gameWindow.UpdateState(stateView)
		case connection_messages.GAME_STATE_INFO:
			var gameStateInfo connection_messages.GameStateInfo
			if err = json.Unmarshal(msg, &gameStateInfo); err != nil {
				log.Println("Err parsing StateView")
				continue
			}
			client.gameWindow.SetGameState(gameStateInfo.GameStateValue)
		default:
			log.Println("Unknown message type")
		}
		log.Printf("Current client (You): %v", client)
	}
}

func (client *Client) sendOnReady(readyState bool) {
	readyMsg, err := json.Marshal(connection_messages.NewReadyMessage(readyState, client.id))
	if err != nil {
		log.Println("Err json.Marshal in toggleReady: ", err)
	}
	err = client.conn.WriteMessage(websocket.TextMessage, readyMsg)
	if err != nil {
		log.Println("Err sending ready: ", err)
	}
}
