package network_server

import (
	"encoding/json"
	"errors"
	"fmt"

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
	clientId, err := server.DecodeClientId(actionMsg)
	if err != nil {
		return err
	}

	if actionType != cm.REARRANGE_CARDS &&
		actionType != cm.DRAW_CARD &&
		!server.clients[clientId].drawnCard {

		server.sendWindowMessage(clientId, "You need to draw card first")
		return nil
	}

	switch actionType {
	case cm.DRAW_CARD:
		var actionDrawMessage cm.ActionDrawMessage
		json.Unmarshal(actionMsg, &actionDrawMessage)
		err := server.handleClientDrawCard(
			actionDrawMessage.ClientId,
			actionDrawMessage.DrawSource,
		)
		return err
	case cm.DISCARD_CARD:
		var actionDiscardMessage cm.ActionDiscardMessage
		json.Unmarshal(actionMsg, &actionDiscardMessage)
		err := server.handleClientDiscardCard(
			actionDiscardMessage.ClientId,
			actionDiscardMessage.Card,
		)
		return err
	case cm.INITIAL_MELD:
		var actionInitialMeldMessage cm.ActionInitialMeldMessage
		json.Unmarshal(actionMsg, &actionInitialMeldMessage)
		err := server.handleClientInitialMeld(
			actionInitialMeldMessage.ClientId,
			actionInitialMeldMessage.Sequences,
		)
		return err
	case cm.REARRANGE_CARDS:
		var actionRearrangeCardsMessage cm.ActionRearrangeCardsMessage
		json.Unmarshal(actionMsg, &actionRearrangeCardsMessage)
		err := server.handleRearrangeCards(
			actionRearrangeCardsMessage.ClientId,
			actionRearrangeCardsMessage.Cards,
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

func (server *Server) DecodeClientId(actionMsg []byte) (int, error) {
	var baseMsg struct {
		ClientId int `json:"client_id"`
	}
	if err := json.Unmarshal(actionMsg, &baseMsg); err != nil {
		return -1, err
	}
	return baseMsg.ClientId, nil
}

func (server *Server) handleClientDrawCard(clientId int, drawSource cm.DRAW_TYPE) error {
	if server.clients[clientId].drawnCard {
		err := server.sendWindowMessage(clientId, "You've drawn a card already!")
		return err
	}
	if drawSource == cm.DRAW_FROM_PILE {
		server.table.PlayerDrawCard(clientId)
		server.clients[clientId].drawnCard = true
		server.SendStateViewAll()
		return nil
	}
	if !server.clients[clientId].hasMelded {
		err := server.sendWindowMessage(clientId, "You need to meld first")
		return err
	}

	server.table.PlayerDrawCardFromDiscard(clientId)
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

	server.clients[clientId].drawnCard = false
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

func (server *Server) handleClientInitialMeld(clientId int, sequences [][]*dm.Card) error {
	isPurePresent := false
	sumPoints := 0
	numCards := 0
	for _, sequence := range sequences {
		if !gm.AreBuildingSequence(sequence) {
			err := server.sendWindowMessage(
				clientId,
				"At least one combination doesn't make sequence",
			)
			return err
		}
		if gm.IsPureSequence(sequence) {
			isPurePresent = true
		}
		sumPoints += gm.SequencePoints(sequence)
		numCards += len(sequence)
	}
	if !isPurePresent {
		err := server.sendWindowMessage(clientId, "You need at least one Pure sequence")
		return err
	} else if sumPoints < gm.MIN_POINTS_TO_MELD {
		errMsg := fmt.Sprintf("You need at least %d points", gm.MIN_POINTS_TO_MELD)
		err := server.sendWindowMessage(clientId, errMsg)
		return err
	} else if numCards == len(server.table.Players[clientId].Hand) {
		err := server.sendWindowMessage(clientId, "You need to place last card to discard pile")
		return err
	}
	// TODO: Handle sequences
	server.clients[clientId].hasMelded = true
	return nil
}

func (server *Server) handleRearrangeCards(clientId int, cards []*dm.Card) error {
	server.table.Players[clientId].SetHand(cards)
	return nil
}
