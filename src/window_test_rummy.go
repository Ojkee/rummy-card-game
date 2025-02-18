package main

import (
	"rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
	"rummy-card-game/src/game_logic/player"
	"rummy-card-game/src/window"
)

func main() {
	w := window.NewWindow()
	player := player.NewPlayer(0)
	player.SetHand([]*dm.Card{
		dm.NewCard(dm.CLUBS, dm.Rank(2)),
		dm.NewCard(dm.DIAMONDS, dm.Rank(11)),
		dm.NewCard(dm.ANY, dm.Rank(13)),
		dm.NewCard(dm.SPADES, dm.Rank(5)),
	})
	w.SetGameState(gm.IN_GAME)
	w.UpdateState(
		*connection_messages.NewStateView(0, dm.NewCardQueue(), dm.NewCardQueue(), player, []int{1}),
	)
	w.MainLoop()
}
