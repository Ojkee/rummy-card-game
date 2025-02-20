package network_server

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

func (server *Server) handleClientAction(actionMsg []byte) error {
	actionType, err := server.DecodeActionType(actionMsg)
	if err != nil {
		return err
	}

	switch actionType {
	case cm.DRAW_CARD:
		var actionDrawMessage cm.ActionDrawMessage
		json.Unmarshal(actionMsg, &actionDrawMessage)
		err := server.handleClientDrawCard(actionDrawMessage.ClientId)
		return err
	case cm.DISCARD_CARD:
		var actionDiscardMessage cm.ActionDiscardMessage
		json.Unmarshal(actionMsg, &actionDiscardMessage)
		err := server.handleClientDiscardCard(
			actionDiscardMessage.ClientId,
			actionDiscardMessage.Card,
		)
		return err
	case cm.UNSUPPORTED:
	default:
		return errors.New("Unsupported/Unimplemented Player Action")
	}
	return nil
}

func (server *Server) DecodeActionType(actionMsg []byte) (cm.ACTION_TYPE, error) {
	var baseMsg struct {
		ActionType cm.ACTION_TYPE `json:"action_type"`
	}
	if err := json.Unmarshal(actionMsg, &baseMsg); err != nil {
		return cm.UNSUPPORTED, err
	}
	return baseMsg.ActionType, nil
}

func (server *Server) handleClientDrawCard(clientId int) error {
	if server.clients[clientId].drawnCard {
		err := server.sendWindowMessage(clientId, "You've drawn a card already!")
		return err
	}
	server.table.PlayerDrawCard(clientId)
	server.clients[clientId].drawnCard = true
	server.SendStateViewAll()
	return nil
}

func (server *Server) handleClientDiscardCard(clientId int, card *dm.Card) error {
	err := server.table.PlayerDiscardCard(clientId, card)
	if err != nil {
		return err
	}
	if server.table.IsWinner(clientId) {
		msg, err := cm.NewGameStateInfo(gm.FINISHED).Json()
		if err != nil {
			return err
		}
		for _, client := range server.clients {
			if client == nil {
				continue
			}
			client.conn.WriteMessage(
				websocket.TextMessage,
				msg,
			)
		}
		server.table.SetState(gm.FINISHED)
		return nil
	}

	server.table.NextTurn()
	server.SendStateViewAll()
	return nil
}

func (server *Server) sendWindowMessage(clientId int, textMsg string) error {
	gameWindowText := cm.NewGameWindowText(textMsg)
	msg, err := gameWindowText.Json()
	if err != nil {
		return err
	}
	server.clients[clientId].conn.WriteMessage(websocket.TextMessage, msg)

	return nil
}
