package window

func (window *Window) drawInRound() {
	for _, playerCard := range window.playerCards {
		playerCard.Draw()
	}
}
