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

	if clientId != server.table.GetTurnId() && actionType != cm.REARRANGE_CARDS {
		err := server.sendWindowMessage(clientId, "Not your turn")
		return err
	}

	if actionType != cm.REARRANGE_CARDS &&
		actionType != cm.DRAW_CARD &&
		!server.clients[clientId].drawnCard {

		err := server.sendWindowMessage(clientId, "You need to draw card first")
		return err
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
		var actionInitialMeldMessage cm.ActionMeldMessage
		json.Unmarshal(actionMsg, &actionInitialMeldMessage)
		err := server.handleClientMeld(
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
	case cm.UPDATE_TABLE_SEQUNCE:
		var updateSequenceMessage cm.ActionUpdateTableSequenceMessage
		json.Unmarshal(actionMsg, &updateSequenceMessage)
		err := server.handleUpdateSequnces(
			updateSequenceMessage.ClientId,
			updateSequenceMessage.SequenceId,
			updateSequenceMessage.CardIdx,
			updateSequenceMessage.Card,
		)
		return err
	case cm.ACTION_UNSUPPORTED:
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
		return cm.ACTION_UNSUPPORTED, err
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
		server.table.ManageDrawpile()
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

func (server *Server) handleClientMeld(clientId int, sequences []*cm.SequenceLocked) error {
	isPurePresent := false
	sumPoints := 0
	numCards := 0
	wrongSeqs := make([]*cm.SequenceLocked, 0)
	for _, sequenceLocked := range sequences {
		sequence := sequenceLocked.Cards
		if !gm.AreBuildingSequence(sequence) {
			wrongSeqs = append(wrongSeqs, sequenceLocked)
		}
		if gm.IsPureSequence(sequence) {
			isPurePresent = true
		}
		sumPoints += gm.SequencePoints(sequence)
		numCards += len(sequence)
	}
	if len(wrongSeqs) > 0 {
		err := server.sendWrongCardsInfo(clientId, wrongSeqs)
		return err
	}
	if !isPurePresent && !server.clients[clientId].hasMelded {
		err := server.sendWindowMessage(clientId, "You need at least one Pure sequence")
		return err
	} else if sumPoints < gm.MIN_POINTS_TO_MELD && !server.clients[clientId].hasMelded {
		errMsg := fmt.Sprintf("You need at least %d points", gm.MIN_POINTS_TO_MELD)
		err := server.sendWindowMessage(clientId, errMsg)
		return err
	} else if numCards == len(server.table.Players[clientId].Hand) {
		err := server.sendWindowMessage(clientId, "You need to place last card to discard pile")
		return err
	}
	for _, sequenceLocked := range sequences {
		var sequenceType gm.SEQUENCE_TYPE
		sequence := sequenceLocked.Cards
		if gm.IsPureSequence(sequence) {
			sequenceType = gm.SEQUENCE_PURE
		} else if gm.IsAscendingSequence(sequence) {
			sequenceType = gm.SEQUENCE_ASCENDING
		} else {
			sequenceType = gm.SEQUENCE_SAME_RANK
		}
		server.table.AddNewSequence(sequence, sequenceType)
		server.table.FilterCards(clientId, sequence)
	}
	server.clients[clientId].hasMelded = true
	err := server.SendStateViewAll()
	return err
}

func (server *Server) sendWrongCardsInfo(clientId int, sequences []*cm.SequenceLocked) error {
	err := server.sendWindowMessage(
		clientId,
		"At least one combination doesn't make sequence",
	)
	if err != nil {
		return err
	}
	wrongCardsHilight := cm.NewWrongCardsHighlight(sequences)
	msg, err := wrongCardsHilight.Json()
	if err != nil {
		return err
	}
	server.clients[clientId].conn.WriteMessage(websocket.TextMessage, msg)
	return nil
}

func (server *Server) handleRearrangeCards(clientId int, cards []*dm.Card) error {
	server.table.Players[clientId].SetHand(cards)
	return nil
}

func (server *Server) handleUpdateSequnces(clientId, sequenceId, cardIdx int, card *dm.Card) error {
	if !server.clients[clientId].hasMelded {
		err := server.sendWindowMessage(clientId, "You need to meld first to do that")
		return err
	}
	err := server.table.HandleAvailableSpotInSequence(
		clientId,
		sequenceId,
		cardIdx,
		card,
	)
	if err != nil {
		return nil
	}
	err = server.SendStateViewAll()
	return err
}

func (server *Server) handleClientDebug(debugMsg []byte) error {
	debugType, err := server.DecodeDebugType(debugMsg)
	if err != nil {
		return err
	}
	switch debugType {
	case cm.DEBUG_RESET:
		server.mu.Lock()
		defer server.mu.Unlock()
		server.resetClients()
		server.table.Reset()
		server.table.InitNewGame()
		server.SendStateViewAll()
		break
	case cm.DEBUG_UNSUPPORTED:
	default:
		return errors.New("Unsupported/Unimplemented Debug Message")
	}
	return nil
}

func (server *Server) DecodeDebugType(actionMsg []byte) (cm.DEBUG_MESSAGE_TYPE, error) {
	var baseMsg struct {
		DebugType cm.DEBUG_MESSAGE_TYPE `json:"debug_type"`
	}
	if err := json.Unmarshal(actionMsg, &baseMsg); err != nil {
		return cm.DEBUG_UNSUPPORTED, err
	}
	return baseMsg.DebugType, nil
}
