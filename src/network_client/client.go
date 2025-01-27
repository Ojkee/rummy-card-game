package network_client

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"

	"rummy-card-game/src/connection_messages"
	"rummy-card-game/src/window"
)

type Client struct {
	id         int
	gameWindow window.Window

	lastStateView connection_messages.StateView
}

func NewClient() *Client {
	return &Client{
		id:         -1,
		gameWindow: *window.NewWindow(),
	}
}

func (client *Client) GetId() int {
	return client.id
}

func (client *Client) SetId(id int) {
	client.id = id
}

func (client *Client) GetLastStateView() connection_messages.StateView {
	return client.lastStateView
}

func (client *Client) UpdateStateView(stateView connection_messages.StateView) {
	client.lastStateView = stateView
}

func (client *Client) Connect() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	defer conn.Close()
	log.Println("Connected!")

	go client.readFromServer(conn)
	go client.writeToServer(conn)
	go client.gameWindow.MainLoop()
	select {}
}

func (client *Client) readFromServer(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
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
			client.UpdateStateView(stateView)
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

func (client *Client) writeToServer(conn *websocket.Conn) {
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		message, _ := consoleReader.ReadString('\n')
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Couldn't send message, err: ", err)
			return
		}
	}
}
