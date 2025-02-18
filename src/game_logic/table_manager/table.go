package table_manager

import (
	"fmt"

	"rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
	"rummy-card-game/src/game_logic/player"
)

type Table struct {
	state        gm.GAME_STATE
	MinPlayers   int
	MaxPlayers   int
	turnID       int
	TemplateDeck dm.Deck
	DrawPile     dm.CardQueue
	DiscardPile  dm.CardQueue
	Players      map[int]*player.Player
}

func NewTable(minPlayers, maxPlayers int) *Table {
	return &Table{
		state:        gm.PRE_START,
		MinPlayers:   minPlayers,
		MaxPlayers:   maxPlayers,
		turnID:       0,
		TemplateDeck: *dm.NewDeck(),
		DrawPile:     *dm.NewCardQueue(),
		DiscardPile:  *dm.NewCardQueue(),
		Players:      make(map[int]*player.Player, 0),
	}
}

func (table *Table) InitNewGame() {
	table.shuffleInitDrawPile()
	table.dealCards()
}

func (table *Table) AddNewPlayer(playerId int) {
	table.Players[playerId] = player.NewPlayer(playerId)
}

func (table *Table) RemovePlayer(playerId int) {
	delete(table.Players, playerId)
}

func (table *Table) shuffleInitDrawPile() {
	_drawPile := make([]*dm.Card, 0)
	numDecks := 2
	for range numDecks {
		for _, card := range *table.TemplateDeck.GetCards() {
			_drawPile = append(_drawPile, &card)
		}
	}
	table.DrawPile.ShuffleExtend(_drawPile)
}

func (table *Table) dealCards() {
	numCardsOnHand := 13
	for playerId := range table.Players {
		hand := make([]*dm.Card, numCardsOnHand)
		for j := range numCardsOnHand {
			hand[j] = table.DrawPile.Pop()
		}
		table.Players[playerId].SetHand(hand)
	}
}

func (table *Table) PrintHands() {
	for _, p := range table.Players {
		fmt.Printf("ID: %d\n", p.Id)
		for _, card := range p.Hand {
			fmt.Printf("%s ", card)
		}
		fmt.Println()
	}
}

func (table *Table) JsonPlayerStateView(playerId int) ([]byte, error) {
	sv := connection_messages.NewStateView(
		table.turnID,
		&table.DrawPile,
		&table.DiscardPile,
		table.Players[playerId],
		[]int{1},
	)
	return sv.Json()
}

func (table *Table) GetTurnId() int {
	return table.turnID
}

func (table *Table) NumPlayers() int {
	return len(table.Players)
}

func (table *Table) NextTurn() {
	table.turnID = (table.turnID + 1) % table.NumPlayers()
}

func (table *Table) GetState() gm.GAME_STATE {
	return table.state
}

func (table *Table) SetState(state gm.GAME_STATE) {
	table.state = state
}

func (table *Table) CanPlayerJoin() bool {
	return table.NumPlayers() < table.MaxPlayers
}

func (table *Table) PlayerDrawCard(playerId int) {
	newCard := table.DrawPile.PopBack()
	table.Players[playerId].DrawCard(newCard)
}
