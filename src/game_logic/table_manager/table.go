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
	NumPlayers   int
	MaxPlayers   int
	turnID       int
	TemplateDeck dm.Deck
	DrawPile     dm.CardQueue
	DiscardPile  dm.CardQueue
	Players      []*player.Player
}

func NewTable(maxPlayers int) *Table {
	_template_deck := *dm.NewDeck()
	return &Table{
		state:        gm.PRE_START,
		NumPlayers:   0,
		MaxPlayers:   maxPlayers,
		turnID:       0,
		TemplateDeck: _template_deck,
		DrawPile:     *dm.NewCardQueue(),
		DiscardPile:  *dm.NewCardQueue(),
		Players:      make([]*player.Player, 0),
	}
}

func (table *Table) InitNewGame() {
	table.shuffleInitDrawPile()
	table.dealCards()
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
	for i := range table.NumPlayers {
		hand := make([]*dm.Card, numCardsOnHand)
		for j := range numCardsOnHand {
			hand[j] = table.DrawPile.Pop()
		}
		table.Players[i].SetHand(hand)
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

func (table *Table) JsonCurrentPlayerStateView() ([]byte, error) {
	sv := connection_messages.NewStateView(
		&table.DrawPile,
		&table.DiscardPile,
		table.Players[table.turnID],
		[]int{1},
	)
	return sv.Json()
}

func (table *Table) GetTurnId() int {
	return table.turnID
}

func (table *Table) NextTurn() {
	table.turnID = (table.turnID + 1) % table.NumPlayers
}

func (table *Table) GetState() gm.GAME_STATE {
	return table.state
}

func (table *Table) SetState(state gm.GAME_STATE) {
	table.state = state
}

func (table *Table) CanPlayerJoin() bool {
	return table.NumPlayers < table.MaxPlayers
}
