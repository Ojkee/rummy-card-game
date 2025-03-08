package table_manager

import (
	"errors"
	"fmt"
	"slices"

	"rummy-card-game/src/connection_messages"
	df "rummy-card-game/src/debug_functools"
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

func (table *Table) Reset() {
	table.turnIdx = 0
	table.TemplateDeck = *dm.NewDeck()
	table.DrawPile = *dm.NewCardQueue()
	table.DiscardPile = *dm.NewCardQueue()
	table.sequences = nil
	table.sequences = make([]gm.Sequence, 0)
	table.jokerImitations = nil
	table.jokerImitations = make(map[int][]gm.JokerImitation)
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
	if df.DEBUG_MODES[df.MELD_HAND_START] {
		skipped := _drawPile[df.SKIP_MELD_HAND_CARDS:]
		n := len(_drawPile)
		_drawPile = append(_drawPile[:n-df.SKIP_MELD_HAND_CARDS], skipped...)
		sameRankDebug := []*dm.Card{
			dm.NewCard(dm.SPADES, dm.FOUR),
			dm.NewCard(dm.HEARTS, dm.FOUR),
			dm.NewCard(dm.CLUBS, dm.FOUR),
			dm.NewCard(dm.DIAMONDS, dm.FOUR),
		}
		sameRankDebug = append(sameRankDebug, _drawPile[df.SKIP_MELD_HAND_CARDS:]...)
		table.DrawPile.Extend(sameRankDebug)
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
		jokImitations := table.getImitationsFromSameRank(cards)
		table.jokerImitations[seqId] = jokImitations
		newSequence := *gm.NewSequence(seqId, cards, sequenceType, jokImitations)
		table.sequences = append(table.sequences, newSequence)
	}
}

// Assume it is propper ascending sequence
func (table *Table) sortAscendingSequence(cards []*dm.Card) ([]*dm.Card, []gm.JokerImitation) {
	jokerPositions := getJokerPositions(cards)
	numJokers := len(jokerPositions)
	sortedCards := gm.SortByRank(cards)
	jokerImitations := make([]gm.JokerImitation, 0)
	if numJokers == 0 {
		return sortedCards, jokerImitations
	}

	numFirstJokers := getNumJokerFirstPositions(jokerPositions)
	numCanFitBefore := int(sortedCards[0].Rank) + 1 // iota starts from: TWO = 0

	sortFilterCards := filterJoks(sortedCards)
	fillCards := insertJokers(sortFilterCards)

	allJoksLeft := numJokers - (len(fillCards) - len(sortFilterCards))
	numBegin := slices.Min([]int{allJoksLeft, numFirstJokers, numCanFitBefore})

	filledBeginGapsCards := appendJoksBeggining(fillCards, numBegin)
	allJoksLeft -= numBegin
	allFilledCards := appendLastJoks(filledBeginGapsCards, allJoksLeft)
	jokerImitations = getImitationsFromSorted(allFilledCards)
	return allFilledCards, jokerImitations
}

func filterJoks(cards []*dm.Card) []*dm.Card {
	retVal := make([]*dm.Card, 0)
	for _, card := range cards {
		if card.Rank != dm.JOKER {
			retVal = append(retVal, card)
		}
	}
	return retVal
}

func insertJokers(filteredCards []*dm.Card) []*dm.Card {
	nextRank := gm.NextRank(filteredCards[0].Rank, true)
	for i := 1; i < len(filteredCards); i++ {
		if filteredCards[i].Rank != *nextRank {
			filteredCards = insert(filteredCards, dm.NewCard(dm.ANY, dm.JOKER), i)
		}
		nextRank = gm.NextRank(*nextRank, false)
	}
	return filteredCards
}

func getJokerPositions(cards []*dm.Card) []int {
	jokerPositions := make([]int, 0)
	for i, card := range cards {
		if card.Rank == dm.JOKER {
			jokerPositions = append(jokerPositions, i)
		}
	}
	return jokerPositions
}

func getNumJokerFirstPositions(jokerPositions []int) int {
	for i, val := range jokerPositions {
		if i != val {
			return i
		}
	}
	return len(jokerPositions)
}

func prepend[T *dm.Card | gm.JokerImitation](slice []T, v T) []T {
	return append([]T{v}, slice...)
}

func insert[T *dm.Card | gm.JokerImitation](slice []T, v T, idx int) []T {
	return append(slice[:idx], prepend(slice[idx:], v)...)
}

func appendJoksBeggining(cards []*dm.Card, numJoks int) []*dm.Card {
	for range numJoks {
		cards = prepend(cards, dm.NewCard(dm.ANY, dm.JOKER))
	}
	return cards
}

func getImitationsFromSorted(cards []*dm.Card) []gm.JokerImitation {
	imitations := make([]gm.JokerImitation, 0)
	lowestCard, lowestIdx := findLowestNonJokCardIdx(cards)
	lowestRank := &lowestCard.Rank
	for i := lowestIdx - 1; i >= 0; i-- {
		lowestRank = gm.PrevRank(*lowestRank, true)
		imitations = prepend(imitations,
			*gm.NewJokerImitation(i, dm.NewCard(lowestCard.Suit, *lowestRank)),
		)
	}
	higherRank := &lowestCard.Rank
	for i := lowestIdx + 1; i < len(cards); i++ {
		higherRank = gm.NextRank(*higherRank, false)
		if cards[i].Rank == dm.JOKER {
			imitations = append(
				imitations,
				*gm.NewJokerImitation(i, dm.NewCard(lowestCard.Suit, *higherRank)),
			)
		}
	}
	return imitations
}

func findLowestNonJokCardIdx(cards []*dm.Card) (*dm.Card, int) {
	for i, card := range cards {
		if card.Rank != dm.JOKER {
			return card, i
		}
	}
	return nil, -1
}

func appendLastJoks(cards []*dm.Card, numJoks int) []*dm.Card {
	lastRank := cards[0].Rank
	for _, card := range cards {
		if card.Rank != dm.JOKER {
			lastRank = card.Rank
		}
	}
	numCanFitAfter := int(dm.ACE) - int(lastRank)
	numFitAfter := min(numCanFitAfter, numJoks)
	for range numFitAfter {
		cards = append(cards, dm.NewCard(dm.ANY, dm.JOKER))
	}
	return appendJoksBeggining(cards, numJoks-numFitAfter)
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
	table.updateJokerImitation(sequenceId)
	return nil
}

func (table *Table) getCardIdFromJokImitations(seqId int, card *dm.Card) (int, error) {
	jokImits, ok := table.jokerImitations[seqId]
	if !ok {
		return -1, errors.New("Sequence to update not found")
	}
	for _, jokImit := range jokImits {
		if *jokImit.Card == *card || (jokImit.CardAlt != nil && *jokImit.CardAlt == *card) {
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

func (table *Table) updateJokerImitation(seqId int) {
	if table.sequences[seqId].Type == gm.SEQUENCE_SAME_RANK {
		cards := table.sequences[seqId].TableCards
		fmt.Println(cards)
		imitations := table.getImitationsFromSameRank(cards)
		table.sequences[seqId].JokerImitations = imitations
		table.jokerImitations[seqId] = imitations
	} else {
		oldSeqCards := table.sequences[seqId].TableCards
		newSeqCards, imitations := table.sortAscendingSequence(oldSeqCards)
		table.sequences[seqId].TableCards = newSeqCards
		table.sequences[seqId].JokerImitations = imitations
		table.jokerImitations[seqId] = imitations
	}
}

func (table *Table) getImitationsFromSameRank(cards []*dm.Card) []gm.JokerImitation {
	jokIdxs := getJokerPositions(cards)
	if len(jokIdxs) == 0 {
		return []gm.JokerImitation{}
	}
	targetRank := table.getRankFromSameRankSequence(cards)
	unusedSuits := table.getUnusedSuits(cards)

	if len(jokIdxs) == 1 && len(unusedSuits) == 2 {
		// 3 CARDS; 1 JOK
		imitCard := dm.NewCard(unusedSuits[0], targetRank)
		imitCardAlt := dm.NewCard(unusedSuits[1], targetRank)
		jokImit := gm.NewJokerImitation(jokIdxs[0], imitCard)
		jokImit.SetCardAlt(imitCardAlt)
		return []gm.JokerImitation{*jokImit}
	} else if len(jokIdxs) == 1 && len(unusedSuits) == 1 {
		// 4 CARDS; 1 JOK
		imitCard := dm.NewCard(unusedSuits[0], targetRank)
		jokImit := gm.NewJokerImitation(jokIdxs[0], imitCard)
		return []gm.JokerImitation{*jokImit}
	} else if len(jokIdxs) == 2 && len(unusedSuits) == 3 {
		// 3 CARDS; 2 JOK
		cardImit1 := dm.NewCard(unusedSuits[0], targetRank)
		jokImit1 := gm.NewJokerImitation(jokIdxs[0], cardImit1)
		cardImit2 := dm.NewCard(unusedSuits[1], targetRank)
		cardImit2Alt := dm.NewCard(unusedSuits[2], targetRank)
		jokImit2 := gm.NewJokerImitation(jokIdxs[1], cardImit2)
		jokImit2.SetCardAlt(cardImit2Alt)
		return []gm.JokerImitation{*jokImit1, *jokImit2}
	}
	// 4 CARDS; 2 JOKS
	cardImit1 := dm.NewCard(unusedSuits[0], targetRank)
	cardImit2 := dm.NewCard(unusedSuits[1], targetRank)
	JokImit1 := gm.NewJokerImitation(jokIdxs[0], cardImit1)
	JokImit2 := gm.NewJokerImitation(jokIdxs[1], cardImit2)
	return []gm.JokerImitation{*JokImit1, *JokImit2}
}

func (table *Table) getRankFromSameRankSequence(cards []*dm.Card) dm.Rank {
	for _, card := range cards {
		if card.Rank != dm.JOKER {
			return card.Rank
		}
	}
	return dm.JOKER
}

func (table *Table) getUnusedSuits(cards []*dm.Card) []dm.Suit {
	suitsUsed := map[dm.Suit]bool{
		dm.HEARTS:   false,
		dm.DIAMONDS: false,
		dm.CLUBS:    false,
		dm.SPADES:   false,
	}
	for _, card := range cards {
		if card.Suit != dm.ANY && !suitsUsed[card.Suit] {
			suitsUsed[card.Suit] = true
		}
	}
	unusedSuits := make([]dm.Suit, 0)
	for suit, used := range suitsUsed {
		if !used {
			unusedSuits = append(unusedSuits, suit)
		}
	}
	return unusedSuits
}

func (table *Table) ManageDrawpile() {
	if !table.DrawPile.IsEmpty() {
		return
	}
	discardedCards := table.DiscardPile.LeaveOnlyLast()
	table.DrawPile.ShuffleExtend(discardedCards)
}
