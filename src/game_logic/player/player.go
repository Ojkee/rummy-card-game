package player

import (
	dm "rummy-card-game/src/game_logic/deck_manager"
)

type Player struct {
	Id   int        `json:"id"`
	Name string     `json:"name"`
	Hand []*dm.Card `json:"hand"`
}

func NewPlayer(id int) *Player {
	return &Player{
		Id:   id,
		Hand: make([]*dm.Card, 0),
	}
}

func (p *Player) SetHand(newHand []*dm.Card) {
	p.Hand = newHand
}

func (p *Player) DrawCard(card *dm.Card) {
	p.Hand = append(p.Hand, card)
}
