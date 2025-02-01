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
	gameWindow  window.Window
	id          int
	gameStarted bool
}

func NewClient() *Client {
	return &Client{
		conn:        nil,
		gameWindow:  *window.NewWindow(),
		id:          -1,
		gameStarted: false,
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

	go client.readFromServer()
	go client.writeToServer()
	go client.gameWindow.MainLoop()

	select {
	case <-client.gameWindow.CloseListener():
		log.Println("Disconnecting!")
	case <-signalChannel:
		client.gameWindow.Stop()
	}
}

func (client *Client) readFromServer() {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading: ", err)
		}

		var messageType connection_messages.MESSAGE_TYPE
		messageType, err = client.decodeMessageType(msg)

		switch messageType {
		case connection_messages.ID_INFO:
			var idInfo connection_messages.IdInfo
			if err = json.Unmarshal(msg, &idInfo); err != nil {
				log.Println("Err parsing Id")
				continue
			}
			client.SetId(idInfo.Id)
			break
		case connection_messages.STATE_VIEW:
			var stateView connection_messages.StateView
			if err = json.Unmarshal(msg, &stateView); err != nil {
				log.Println("Err parsing StateView")
				continue
			}
			client.gameWindow.UpdateState(stateView)
			break
		default:
			break
		}
		log.Printf("Current client (You): %v", client)
	}
}

func (client *Client) decodeMessageType(msg []byte) (connection_messages.MESSAGE_TYPE, error) {
	var messageDecoded map[string]json.RawMessage
	var err error
	err = json.Unmarshal(msg, &messageDecoded)
	if err != nil {
		return connection_messages.UNKNOWN, err
	}
	var messageType connection_messages.MESSAGE_TYPE
	err = json.Unmarshal(messageDecoded["message_type"], &messageType)
	if err != nil {
		return connection_messages.UNKNOWN, err
	}
	return messageType, nil
}

func (client *Client) writeToServer() {
	for {
	}
}

func (client *Client) sendOnReady(readyState bool) {
	readyMsg, err := json.Marshal(connection_messages.NewReadyMessage(readyState))
	if err != nil {
		log.Println("Err json.Marshal in toggleReady: ", err)
	}
	err = client.conn.WriteMessage(websocket.TextMessage, readyMsg)
	if err != nil {
		log.Println("Err sending ready: ", err)
	}
}
