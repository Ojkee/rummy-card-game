package table_manager

import (
	"errors"
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
	turnId       int
	turnIdx      int
	TemplateDeck dm.Deck
	DrawPile     dm.CardQueue
	DiscardPile  dm.CardQueue
	Players      map[int]*player.Player
	playerIds    []int

	sequences       []gm.Sequence
	jokerImitations map[int][]gm.JokerImitation // SEQ_ID: JOK_IMITS
}

func NewTable(minPlayers, maxPlayers int) *Table {
	return &Table{
		state:        gm.PRE_START,
		MinPlayers:   minPlayers,
		MaxPlayers:   maxPlayers,
		turnId:       -1,
		turnIdx:      0,
		TemplateDeck: *dm.NewDeck(),
		DrawPile:     *dm.NewCardQueue(),
		DiscardPile:  *dm.NewCardQueue(),
		Players:      make(map[int]*player.Player, 0),
		playerIds:    make([]int, 0),

		sequences:       make([]gm.Sequence, 0),
		jokerImitations: make(map[int][]gm.JokerImitation),
	}
}

func (table *Table) InitNewGame() {
	table.shuffleInitDrawPile()
	table.dealCards()
	table.initPlayersIds()
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
	if DEBUG_MODES[MELD_HAND_START] {
		skipped := _drawPile[SKIP_MELD_HAND_CARDS:]
		n := len(_drawPile)
		_drawPile = append(_drawPile[:n-SKIP_MELD_HAND_CARDS], skipped...)
		table.DrawPile.Extend(_drawPile[SKIP_MELD_HAND_CARDS:])
	} else {
		table.DrawPile.ShuffleExtend(_drawPile)
	}
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
		table.turnId,
		&table.DrawPile,
		&table.DiscardPile,
		table.Players[playerId],
		[]int{1},
		table.sequences,
	)
	return sv.Json()
}

func (table *Table) GetTurnId() int {
	return table.turnId
}

func (table *Table) NumPlayers() int {
	return len(table.Players)
}

func (table *Table) initPlayersIds() {
	table.turnIdx = 0
	table.playerIds = make([]int, 0)
	for playerId := range table.Players {
		table.playerIds = append(table.playerIds, playerId)
	}
	if len(table.playerIds) == 0 {
		return
	}
	table.turnId = table.playerIds[table.turnIdx]
}

func (table *Table) NextTurn() {
	table.turnIdx = (table.turnIdx + 1) % table.NumPlayers()
	table.turnId = table.playerIds[table.turnIdx]
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

func (table *Table) PlayerDrawCardFromDiscard(playerId int) {
	newCard := table.DiscardPile.PopBack()
	table.Players[playerId].DrawCard(newCard)
}

func (table *Table) PlayerDiscardCard(playerId int, discardedCard *dm.Card) error {
	before := len(table.Players[playerId].Hand)
	resultHand := make([]*dm.Card, 0)
	discarded := false
	for _, card := range table.Players[playerId].Hand {
		if *card != *discardedCard || discarded {
			resultHand = append(resultHand, card)
		} else if *card == *discardedCard {
			discarded = true
		}
	}
	after := len(resultHand)
	if before == after {
		return errors.New("No card removed")
	}
	table.Players[playerId].SetHand(resultHand)
	table.DiscardPile.Push(discardedCard)
	return nil
}

func (table *Table) IsWinner(playerId int) bool {
	return len(table.Players[playerId].Hand) == 0
}

func (table *Table) AddNewSequence(cards []*dm.Card, sequenceType gm.SEQUENCE_TYPE) {
	seqId := len(table.sequences)
	if sequenceType != gm.SEQUENCE_SAME_RANK && gm.ContainsJoker(cards) {
		sortedCards, jokImitations := table.sortAscendingSequence(cards)
		table.jokerImitations[seqId] = jokImitations
		newSequence := *gm.NewSequence(seqId, sortedCards, sequenceType, jokImitations)
		table.sequences = append(table.sequences, newSequence)
	} else {
		table.sequences = append(table.sequences, *gm.NewSequence(seqId, cards, sequenceType, []gm.JokerImitation{}))
	}
}

func (table *Table) sortAscendingSequence(cards []*dm.Card) ([]*dm.Card, []gm.JokerImitation) {
	jokerImitations := make([]gm.JokerImitation, 0)
	sortedCards := gm.SortByRank(cards)
	nextRank := gm.NextRank(sortedCards[0].Rank, true)
	n := len(sortedCards)
	var suit dm.Suit
	for _, card := range sortedCards {
		if card.Suit != dm.ANY {
			suit = card.Suit
			break
		}
	}
	for i := 1; i < n; i++ {
		if sortedCards[i].Rank == dm.JOKER {
			// TODO: current card add imitation
			return sortedCards, jokerImitations
		}
		if sortedCards[i].Rank != *nextRank {
			jokFromEnd := sortedCards[n-1]
			for j := n - 1; j > i; j-- {
				sortedCards[j] = sortedCards[j-1]
			}
			sortedCards[i] = jokFromEnd
			imitatedCard := dm.NewCard(suit, *nextRank)
			jokerImitation := *gm.NewJokerImitation(i, imitatedCard)
			jokerImitations = append(jokerImitations, jokerImitation)
		}
		nextRank = gm.NextRank(*nextRank, false)
	}
	return sortedCards, jokerImitations
}

func (table *Table) FilterCards(playerId int, cards []*dm.Card) {
	for _, card := range cards {
		table.filterCard(playerId, card)
	}
}

func (table *Table) filterCard(playerId int, filterCard *dm.Card) {
	resultHand := make([]*dm.Card, 0)
	discarded := false
	for _, card := range table.Players[playerId].Hand {
		if *card != *filterCard || discarded {
			resultHand = append(resultHand, card)
		} else if *card == *filterCard {
			discarded = true
		}
	}
	table.Players[playerId].SetHand(resultHand)
}

func (table *Table) HandleAvailableSpotInSequence(
	playerId, sequenceId, cardIdx int,
	card *dm.Card,
) error {
	if cardIdx < 0 {
		table.appendBeginSequence(card, sequenceId)
	} else if cardIdx > len(table.sequences[sequenceId].TableCards) {
		table.sequences[sequenceId].TableCards = append(
			table.sequences[sequenceId].TableCards,
			card,
		)
	} else {
		cardJokIdx, err := table.getCardIdFromJokImitations(sequenceId, card)
		if err != nil {
			return err
		}
		replaceCard := table.sequences[sequenceId].TableCards[cardJokIdx]
		table.Players[playerId].Hand = append(table.Players[playerId].Hand, replaceCard)
		table.sequences[sequenceId].TableCards[cardJokIdx] = card
		table.filterJokImitation(sequenceId, cardJokIdx)
	}

	table.filterCard(playerId, card)
	return nil
}

func (table *Table) getCardIdFromJokImitations(seqId int, card *dm.Card) (int, error) {
	jokImits, ok := table.jokerImitations[seqId]
	if !ok {
		return -1, errors.New("Sequence to update not found")
	}
	for _, jokImit := range jokImits {
		if *jokImit.Card == *card {
			return jokImit.Idx, nil
		}
	}
	errMsg := fmt.Sprintf("Card: %v not found in imitations", card)
	return -1, errors.New(errMsg)
}

func (table *Table) appendBeginSequence(card *dm.Card, sequenceId int) {
	newSeq := []*dm.Card{card}
	newSeq = append(newSeq, table.sequences[sequenceId].TableCards...)
	table.sequences[sequenceId].TableCards = newSeq
}

func (table *Table) filterJokImitation(seqId, cardId int) {
	newJokImits := make([]gm.JokerImitation, 0)
	for _, jokImit := range table.jokerImitations[seqId] {
		if jokImit.Idx != cardId {
			newJokImits = append(newJokImits, jokImit)
		}
	}
	table.jokerImitations[seqId] = newJokImits
}
