package network_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"rummy-card-game/src/connection_messages"
	"rummy-card-game/src/game_logic/game_manager"
	tm "rummy-card-game/src/game_logic/table_manager"
)

type Server struct {
	upgrader        *websocket.Upgrader
	mu              sync.Mutex
	roomId          int
	gameStarted     bool
	table           *tm.Table
	connections     map[int]*websocket.Conn
	readinessStates map[int]bool
}

func NewServer(maxPlayers int) *Server {
	_connections := make(map[int]*websocket.Conn)
	for i := range maxPlayers {
		_connections[i] = nil
	}
	_readinessStates := make(map[int]bool)
	for i := range maxPlayers {
		_readinessStates[i] = false
	}
	return &Server{
		mu:              sync.Mutex{},
		roomId:          0,
		gameStarted:     false,
		table:           tm.NewTable(maxPlayers),
		connections:     _connections,
		readinessStates: _readinessStates,
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
	server.table.AddNewPlayer(playerId)

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
				server.table.NumPlayers,
				server.table.MaxPlayers,
			)
			return
		}

		messageType, err := connection_messages.DecodeMessageType(msg)
		if err != nil {
			log.Println("Error decoding message type:", err)
			continue
		}

		switch messageType {
		case connection_messages.PLAYER_ACTION:
			continue
		case connection_messages.PLAYER_READY:
			var readyMessage connection_messages.ReadyMessage
			json.Unmarshal(msg, &readyMessage)
			server.manageReadinessStates(readyMessage.ClientId, readyMessage.IsReady)
		}
	}
}

func (server *Server) manageReadinessStates(clientId int, state bool) {
	server.readinessStates[clientId] = state
	if server.allReady() {
		log.Println("ALL READY")
		for _, conn := range server.connections {
			msg, err := connection_messages.NewGameStateInfo(game_manager.IN_GAME).Json()
			if err != nil {
				return
			}
			conn.WriteMessage(
				websocket.TextMessage,
				msg,
			)
		}
	}
}

func (server *Server) addPlayerConnection(conn *websocket.Conn) (int, bool) {
	server.mu.Lock()
	defer server.mu.Unlock()
	playerId, ok := server.GetNextAvailablePlayerId()
	if !ok {
		return -1, false
	}
	server.table.NumPlayers += 1
	log.Printf(
		"Player joined, Id: %d\n\tPlayers: %d/%d\n",
		playerId,
		server.table.NumPlayers,
		server.table.MaxPlayers,
	)
	server.connections[playerId] = conn
	return playerId, true
}

func (server *Server) deletePlayerConnection(playerId int) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.connections[playerId] = nil
	server.table.NumPlayers -= 1
}

func (server *Server) GetNextAvailablePlayerId() (int, bool) {
	for key, val := range server.connections {
		if val == nil {
			return key, true
		}
	}
	return -1, false
}

func (server *Server) SendTableJson(playerId int) {
	server.mu.Lock()
	defer server.mu.Unlock()

	stateView, err := server.table.JsonPlayerStateView(playerId)
	if err != nil {
		log.Println("Err: ", err)
		return
	}
	err = server.connections[playerId].WriteMessage(websocket.TextMessage, stateView)
	if err != nil {
		log.Println("Err: ", err)
		return
	}
}

func (server *Server) SendIdJson(conn *websocket.Conn, id int) {
	server.mu.Lock()
	defer server.mu.Unlock()

	idInfo := connection_messages.NewIdInfo(id)
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
	playerCounter := 0
	for playerId, val := range server.readinessStates {
		if conn, ok := server.connections[playerId]; ok && conn != nil && !val {
			return false
		} else if ok && conn != nil && val {
			playerCounter++
		}
	}
	return playerCounter == server.table.MaxPlayers
}
