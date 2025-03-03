package network_client

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	cm "rummy-card-game/src/connection_messages"
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
	client.gameWindow.SetClientId(id)
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
	client.gameWindow.SetActionMessageCallback(client.sendActionMessage)
	client.gameWindow.SetDebugMessageCallback(client.sendDebugMessage)

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

		var messageType cm.MESSAGE_TYPE
		messageType, err = cm.DecodeMessageType(msg)
		if err != nil {
			log.Println("Error decoding message type:", err)
			continue
		}

		switch messageType {
		case cm.ID_INFO:
			var idInfo cm.IdInfo
			if err = json.Unmarshal(msg, &idInfo); err != nil {
				log.Println("Err parsing Id")
				continue
			}
			client.SetId(idInfo.Id)
		case cm.STATE_VIEW:
			var stateView cm.StateView
			if err = json.Unmarshal(msg, &stateView); err != nil {
				log.Println("Err parsing StateView")
				continue
			}
			client.gameWindow.UpdateState(stateView)
		case cm.GAME_STATE_INFO:
			var gameStateInfo cm.GameStateInfo
			if err = json.Unmarshal(msg, &gameStateInfo); err != nil {
				log.Println("Err parsing StateView")
				continue
			}
			client.gameWindow.SetGameState(gameStateInfo.GameStateValue)
		case cm.GAME_WINDOW_TEXT:
			var gameWindowText cm.GameWindowText
			if err = json.Unmarshal(msg, &gameWindowText); err != nil {
				log.Println("Err parsing GameWindowText")
				continue
			}
			client.gameWindow.PlaceText(gameWindowText.Value)
		case cm.WRONG_CARDS_HIGHLIGHT:
			var wrongCardsHighlight cm.WrongCardsHighlight
			if err = json.Unmarshal(msg, &wrongCardsHighlight); err != nil {
				log.Println("Err parsing WrongCardsHighlight")
				continue
			}
			client.gameWindow.PlaceWrongCardsHighlight(wrongCardsHighlight.SeqLocked)
		default:
			log.Println("Unknown message type")
		}
	}
}

func (client *Client) sendOnReady(readyState bool) {
	readyMsg, err := json.Marshal(cm.NewReadyMessage(readyState, client.id))
	if err != nil {
		log.Println("Err json.Marshal in toggleReady: ", err)
		return
	}
	err = client.conn.WriteMessage(websocket.TextMessage, readyMsg)
	if err != nil {
		log.Println("Err sending ready: ", err)
		return
	}
}

func (client *Client) sendActionMessage(actionMsg cm.ActionMessage) {
	msg, err := actionMsg.Json()
	if err != nil {
		log.Println("Err action message json")
		return
	}
	client.conn.WriteMessage(
		websocket.TextMessage,
		msg,
	)
}

func (client *Client) sendDebugMessage(debugMsg cm.DebugMessage) {
	msg, err := debugMsg.Json()
	if err != nil {
		log.Println("Err debug message json")
		return
	}
	client.conn.WriteMessage(
		websocket.TextMessage,
		msg,
	)
}
