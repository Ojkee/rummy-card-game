package network_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	cm "rummy-card-game/src/connection_messages"
	gm "rummy-card-game/src/game_logic/game_manager"
	tm "rummy-card-game/src/game_logic/table_manager"
)

type Server struct {
	debugMode DEBUG_MODE

	upgrader    *websocket.Upgrader
	mu          sync.Mutex
	roomId      int
	gameStarted bool
	table       *tm.Table
	clients     map[int]*ConnectedClient
}

func NewServer(minPlayers, maxPlayers int) *Server {
	_clients := make(map[int]*ConnectedClient)
	for i := range maxPlayers {
		_clients[i] = nil
	}
	return &Server{
		debugMode: NO_DEBUG,

		mu:          sync.Mutex{},
		roomId:      0,
		gameStarted: false,
		table:       tm.NewTable(minPlayers, maxPlayers),
		clients:     _clients,
	}
}

func (server *Server) Init(port string) {
	server.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return server.table.CanPlayerJoin() },
	}
	http.HandleFunc("/ws", server.HandleConnection)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("Error starting: ", err)
	}
}

func (server *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket: ", err)
		return
	}
	go server.handleClient(conn)
}

func (server *Server) handleClient(conn *websocket.Conn) {
	defer conn.Close()
	if !server.table.CanPlayerJoin() {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Room is full")); err != nil {
			log.Println("Error writing message: ", err)
		}
		return
	}
	playerId, ok := server.addPlayerConnection(conn)
	if !ok {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Error giving Id")); err != nil {
			log.Println("Error writing message: ", err)
		}
		return
	}

	server.SendIdJson(conn, playerId)
	go server.readFromClient(conn, playerId)

	select {}
}

func (server *Server) readFromClient(conn *websocket.Conn, playerId int) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			server.deletePlayerConnection(playerId)
			log.Printf(
				"Player disconnected, Id: %d\n\tPlayers: %d/%d\n",
				playerId,
				server.table.NumPlayers(),
				server.table.MaxPlayers,
			)
			return
		}

		messageType, err := cm.DecodeMessageType(msg)
		if err != nil {
			log.Println("Error decoding message type:", err)
			continue
		}

		switch messageType {
		case cm.PLAYER_ACTION:
			decodedClientId, err := cm.DecodeMessageClientId(msg)
			if decodedClientId != server.table.GetTurnId() {
				server.sendWindowMessage(decodedClientId, "Not your turn")
				break
			}
			err = server.handleClientAction(msg)
			if err != nil {
				log.Println(err)
				break
			}
		case cm.PLAYER_READY:
			var readyMessage cm.ReadyMessage
			json.Unmarshal(msg, &readyMessage)
			server.manageReadinessStates(readyMessage.ClientId, readyMessage.IsReady)
		default:
			continue
		}
	}
}

func (server *Server) manageReadinessStates(clientId int, state bool) {
	server.clients[clientId].isReady = state
	if server.allReady() {
		log.Println("ALL READY")
		for _, client := range server.clients {
			if client == nil {
				continue
			}
			msg, err := cm.NewGameStateInfo(gm.IN_GAME).Json()
			if err != nil {
				continue
			}
			client.conn.WriteMessage(
				websocket.TextMessage,
				msg,
			)
		}
		server.table.InitNewGame()
		err := server.SendStateViewAll()
		if err != nil {
			log.Println(err)
			return
		}
		server.table.SetState(gm.IN_GAME)
	}
}

func (server *Server) addPlayerConnection(conn *websocket.Conn) (int, bool) {
	server.mu.Lock()
	defer server.mu.Unlock()
	playerId, ok := server.GetNextAvailablePlayerId()
	if !ok {
		return -1, false
	}
	server.table.AddNewPlayer(playerId)
	log.Printf(
		"Player joined, Id: %d\n\tPlayers: %d/%d\n",
		playerId,
		server.table.NumPlayers(),
		server.table.MaxPlayers,
	)
	server.clients[playerId] = NewConnectedClient(conn)
	return playerId, true
}

func (server *Server) deletePlayerConnection(playerId int) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.clients[playerId] = nil
	server.table.RemovePlayer(playerId)
}

func (server *Server) GetNextAvailablePlayerId() (int, bool) {
	for playerId, client := range server.clients {
		if client == nil {
			return playerId, true
		}
	}
	return -1, false
}

func (server *Server) SendIdJson(conn *websocket.Conn, id int) {
	server.mu.Lock()
	defer server.mu.Unlock()

	idInfo := cm.NewIdInfo(id)
	idInfoJson, err := idInfo.Json()
	if err != nil {
		log.Println("Err: ", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, idInfoJson)
	if err != nil {
		log.Println("Err: ", err)
		return
	}
}

func (server *Server) allReady() bool {
	readyCounter := 0
	for _, client := range server.clients {
		if client == nil {
			continue
		}
		if !client.isReady {
			return false
		}
		readyCounter++
	}
	return readyCounter >= server.table.MinPlayers
}

func (server *Server) SendStateViewAll() error {
	for playerId, client := range server.clients {
		if client == nil || client.conn == nil {
			continue
		}
		sv, err := server.table.JsonPlayerStateView(playerId)
		if err != nil {
			return err
		}
		err = client.conn.WriteMessage(websocket.TextMessage, sv)
		if err != nil {
			return err
		}
	}
	return nil
}
