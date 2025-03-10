package game_manager

import (
	"sort"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type GAME_STATE int

const (
	PRE_CONNECT GAME_STATE = iota
	PRE_START
	IN_GAME
	FINISHED
)

const MIN_POINTS_TO_MELD = 51

func AreBuildingSequence(cards []*dm.Card) bool {
	if allJoks(cards) || len(cards) < 3 {
		return false
	}
	cardsCopy := make([]*dm.Card, len(cards))
	copy(cardsCopy, cards)
	return IsAscendingSequence(cardsCopy) || IsSameRankSequence(cardsCopy)
}

func allJoks(cards []*dm.Card) bool {
	for _, card := range cards {
		if card.Rank != dm.JOKER {
			return false
		}
	}
	return true
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
	cardsCopy := make([]*dm.Card, len(cards))
	copy(cardsCopy, cards)
	sortedCards := SortByRank(cardsCopy)
	if !areSameSuit(sortedCards) {
		return false
	}
	if len(sortedCards) == 13 {
		return true
	}
	if len(sortedCards) > 13 {
		return false
	}

	return isAscendingNonAce(sortedCards) || isAscendingAce(sortedCards)
}

func isAscendingNonAce(cards []*dm.Card) bool {
	targetRank := NextRank(cards[0].Rank, false)
	if targetRank == nil {
		return false
	}
	usedJoks := 0
	for i := 1; i < len(cards)-usedJoks; {
		if targetRank == nil {
			if cards[i].Rank != dm.JOKER {
				return false
			}
			return true
		}
		if *targetRank != cards[i].Rank {
			if cards[len(cards)-1-usedJoks].Rank != dm.JOKER {
				return false
			} else {
				usedJoks++
			}
		} else {
			i++
		}
		targetRank = NextRank(*targetRank, false)
	}
	return true
}

func isAscendingAce(cards []*dm.Card) bool {
	if !IsRankPresent(dm.ACE, cards) {
		return false
	}
	targetRank := NextRank(dm.ACE, true)
	wasAce := false
	usedJoks := 0
	for i := 0; i < len(cards)-usedJoks; {
		if !wasAce && cards[i].Rank == dm.ACE {
			wasAce = true
			i++
			continue
		} else if wasAce && cards[i].Rank == dm.ACE {
			return false
		} else if *targetRank != cards[i].Rank {
			if cards[len(cards)-1-usedJoks].Rank != dm.JOKER {
				return false
			} else {
				usedJoks++
			}
		}
		targetRank = NextRank(*targetRank, false)
		i++
	}

	return true
}

func IsRankPresent(rank dm.Rank, cards []*dm.Card) bool {
	for _, card := range cards {
		if card.Rank == rank {
			return true
		}
	}
	return false
}

func areSameSuit(cards []*dm.Card) bool {
	suit := cards[0].Suit
	if suit == dm.ANY {
		return false
	}

	for _, card := range cards {
		if card.Suit != suit && card.Suit != dm.ANY {
			return false
		}
	}
	return true
}

func IsPureSequence(cards []*dm.Card) bool {
	cardsCopy := make([]*dm.Card, len(cards))
	copy(cardsCopy, cards)
	if len(cards) < 3 || !IsAscendingSequence(cardsCopy) {
		return false
	}
	nonJokStrek := 0
	for _, card := range cards {
		if card.Rank == dm.JOKER {
			nonJokStrek = 0
		} else {
			nonJokStrek++
		}
		if nonJokStrek >= 3 {
			return true
		}
	}
	return false
}

func SortByRank(cards []*dm.Card) []*dm.Card {
	sort.Slice(cards, func(i int, j int) bool {
		return cards[i].Rank < cards[j].Rank
	})
	return cards
}

func PrevRank(rank dm.Rank, isFirst bool) *dm.Rank {
	if rank == dm.TWO {
		prevRank := dm.ACE
		return &prevRank
	} else if (isFirst && rank == dm.ACE) || rank == dm.JOKER {
		return nil
	}
	prev, _ := dm.RankOfInt(int(rank - 1))
	return &prev
}

func NextRank(rank dm.Rank, isFirst bool) *dm.Rank {
	if rank == dm.ACE && isFirst {
		next := dm.TWO
		return &next
	} else if rank == dm.ACE {
		return nil
	}
	next, _ := dm.RankOfInt(int(rank + 1))
	return &next
}

func SequencePoints(cards []*dm.Card) int {
	sumPoints := 0
	for _, card := range cards {
		sumPoints += card.Rank.Points()
	}
	return sumPoints
}

func ContainsJoker(cards []*dm.Card) bool {
	for _, card := range cards {
		if card.Rank == dm.JOKER {
			return true
		}
	}
	return false
}

func FitSequenceIds(card *dm.Card, sequence *Sequence) []int {
	switch sequence.Type {
	case SEQUENCE_SAME_RANK:
		if (usedSuit(&card.Suit, sequence) || !sameRank(&card.Rank, sequence)) ||
			(len(sequence.TableCards) >= 4 && numJokers(sequence.TableCards) == 0) {
			return []int{}
		} else if len(sequence.TableCards) == 4 {
			return getReplaceIds(card, sequence)
		}
		availableIdxs := []int{-1, len(sequence.TableCards)}
		availableIdxs = append(availableIdxs, getReplaceIds(card, sequence)...)
		return availableIdxs
	case SEQUENCE_PURE, SEQUENCE_ASCENDING:
		if seqSuit := sequence.GetSuitIfAscending(); seqSuit == dm.ANY ||
			(card.Rank != dm.JOKER && seqSuit != card.Suit) {
			return []int{}
		}
		return findInAscending(card, sequence)
	}
	return []int{}
}

func numJokers(cards []*dm.Card) int {
	joks := 0
	for _, card := range cards {
		if card.Rank == dm.JOKER {
			joks++
		}
	}
	return joks
}

func usedSuit(suit *dm.Suit, sequence *Sequence) bool {
	if *suit == dm.ANY {
		return false
	}
	for _, card := range sequence.TableCards {
		if card.Suit == *suit {
			return true
		}
	}
	return false
}

func sameRank(rank *dm.Rank, sequence *Sequence) bool {
	if *rank == dm.JOKER {
		return true
	}
	for _, card := range sequence.TableCards {
		if card.Rank != dm.JOKER && card.Rank != *rank {
			return false
		}
	}
	return true
}

func findInAscending(card *dm.Card, sequence *Sequence) []int {
	retVal := make([]int, 0)
	if canAddBegin(card, sequence) {
		retVal = append(retVal, -1)
	}
	if canAddEnd(card, sequence) {
		retVal = append(retVal, len(sequence.TableCards))
	}
	replaceIds := getReplaceIds(card, sequence)
	retVal = append(retVal, replaceIds...)
	return retVal
}

func canAddBegin(card *dm.Card, sequence *Sequence) bool {
	if card.Rank == dm.JOKER {
		return sequence.TableCards[0].Rank != dm.ACE
	}
	firstCard := sequence.TableCards[0]
	nextRank := NextRank(card.Rank, false)
	if nextRank == nil {
		return false
	}
	if firstCard.Rank != dm.ACE && *nextRank == firstCard.Rank {
		return true
	}
	for _, jokImit := range sequence.JokerImitations {
		if jokImit.Idx == 0 && jokImit.Card.Rank == *nextRank {
			return true
		}
	}
	return false
}

func canAddEnd(card *dm.Card, sequence *Sequence) bool {
	lastIdx := len(sequence.TableCards) - 1
	if card.Rank == dm.JOKER {
		return sequence.TableCards[lastIdx].Rank != dm.ACE
	}
	lastCard := sequence.TableCards[lastIdx]
	if lastCard.Rank != dm.ACE && *NextRank(lastCard.Rank, false) == card.Rank {
		return true
	}
	for _, jokImit := range sequence.JokerImitations {
		if jokImit.Idx == lastIdx && *NextRank(jokImit.Card.Rank, false) == card.Rank {
			return true
		}
	}
	return false
}

func getReplaceIds(card *dm.Card, sequence *Sequence) []int {
	replaceIds := make([]int, 0)
	for _, jokerImitation := range sequence.JokerImitations {
		if *jokerImitation.Card == *card ||
			(jokerImitation.CardAlt != nil && *jokerImitation.CardAlt == *card) {
			return []int{jokerImitation.Idx}
		}
	}
	return replaceIds
}
