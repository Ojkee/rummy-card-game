package game_manager

import (
	"sort"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type GAME_STATE int

const (
	PRE_START GAME_STATE = iota
	IN_GAME
	FINISHED
)

func AreBuildingSequence(cards []*dm.Card) bool {
	if len(cards) < 3 {
		return false
	}
	return IsAscendingSequence(cards) || IsSameRankSequence(cards)
}

func IsSameRankSequence(cards []*dm.Card) bool {
	if len(cards) > 4 {
		return false
	}
	targetRank := cards[0].Rank
	usedSuits := make(map[dm.Suit]bool, 0)
	for _, card := range cards {
		_, isPresentSuit := usedSuits[card.Suit]
		if isPresentSuit || (card.Rank != targetRank && card.Rank != dm.JOKER) {
			return false
		}
		if card.Suit != dm.ANY {
			usedSuits[card.Suit] = true
		}
	}
	return true
}

func IsAscendingSequence(cards []*dm.Card) bool {
	sortedCards := sortByRank(cards)
	targetSuit := sortedCards[0].Suit
	targetRank := nextRank(sortedCards[0].Rank, true)
	if targetRank == nil {
		return false
	}
	usedJokers := 0
	i := 1
	n := len(sortedCards)
	for i < n-usedJokers {
		card := sortedCards[i]
		if card.Rank == dm.JOKER {
			i++
			continue
		}
		if targetRank == nil {
			return false
		}
		if card.Suit != targetSuit && card.Suit != dm.ANY {
			return false
		}
		if *targetRank != card.Rank {
			// Assume that jokers are at the end after sort
			if sortedCards[n-1-usedJokers].Rank == dm.JOKER {
				targetRank = nextRank(*targetRank, false)
				usedJokers++
				continue
			}
			return false
		}
		targetRank = nextRank(*targetRank, false)
		i++
	}
	return true
}

func IsClearSequence(cards []*dm.Card) bool {
	if len(cards) < 3 || !IsAscendingSequence(cards) {
		return false
	}
	for _, card := range cards {
		if card.Rank == dm.JOKER {
			return false
		}
	}
	return true
}

func sortByRank(cards []*dm.Card) []*dm.Card {
	sort.Slice(cards, func(i int, j int) bool {
		return cards[i].Rank < cards[j].Rank
	})
	return cards
}

func nextRank(rank dm.Rank, isFirst bool) *dm.Rank {
	if rank == dm.ACE && !isFirst {
		return nil
	} else if rank == dm.ACE {
		next := dm.TWO
		return &next
	}
	next, _ := dm.RankOfInt(int(rank + 1))
	return &next
}
