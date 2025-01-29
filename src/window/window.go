package window

import (
	"rummy-card-game/src/connection_messages"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type Window struct {
	mu sync.Mutex

	running     bool
	stopChannel chan struct{}

	playerCards       []CardModel
	discardPile       *dm.CardQueue
	lastDiscardedCard *CardModel
}

func NewWindow() *Window {
	return &Window{
		mu: sync.Mutex{},

		running:     true,
		stopChannel: make(chan struct{}),

		playerCards:       make([]CardModel, 0),
		discardPile:       nil,
		lastDiscardedCard: nil,
	}
}

func (window *Window) MainLoop() {
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "Rummy")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	window.initGraphics()

	for !rl.WindowShouldClose() && window.running {
		window.checkEvent()
		window.draw()
	}
	window.Stop()
}

func (window *Window) checkEvent() {
	if rl.IsKeyPressed(rl.KeyQ) {
		window.Stop()
	}
}

func (window *Window) draw() {
	rl.BeginDrawing()

	rl.ClearBackground(COLOR_DARK_GRAY)
	for _, card := range window.playerCards {
		card.Draw()
	}
	rl.EndDrawing()
}

func (window *Window) initGraphics() {
	rl.ImageResize(CLUBS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(DIAMONDS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(HEARTS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(SPADES_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(JOKER_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	RANK_IMGS[dm.CLUBS] = rl.LoadTextureFromImage(CLUBS_IMG)
	RANK_IMGS[dm.DIAMONDS] = rl.LoadTextureFromImage(DIAMONDS_IMG)
	RANK_IMGS[dm.HEARTS] = rl.LoadTextureFromImage(HEARTS_IMG)
	RANK_IMGS[dm.SPADES] = rl.LoadTextureFromImage(SPADES_IMG)
	RANK_IMGS[dm.ANY] = rl.LoadTextureFromImage(JOKER_IMG)
	FONT = rl.LoadFontEx(FONT_PATH, FONT_SIZE, nil, 96)
}

func (window *Window) unloadGraphics() {
	for _, texture := range RANK_IMGS {
		rl.UnloadTexture(texture)
	}
	for _, image := range []*rl.Image{CLUBS_IMG, DIAMONDS_IMG, HEARTS_IMG, SPADES_IMG, JOKER_IMG} {
		rl.UnloadImage(image)
	}
	rl.UnloadFont(FONT)
}

func (window *Window) CloseListener() <-chan struct{} {
	return window.stopChannel
}

func (window *Window) Stop() {
	window.unloadGraphics()
	window.running = false
	close(window.stopChannel)
}

func (window *Window) UpdateState(sv connection_messages.StateView) {
	window.mu.Lock()
	defer window.mu.Unlock()

	window.updatePlayerHand(sv.PlayerEntity.Hand)
	window.discardPile = sv.DiscardPile
	window.updateLastDiscardedCard(window.discardPile.PopBack())
}

func (window *Window) updatePlayerHand(hand []*dm.Card) {
	window.playerCards = make([]CardModel, 0)
	numCards := len(hand)
	offsetX := float32(WINDOW_WIDTH-int32(numCards)*CARD_WIDTH) / 2
	if numCards != 0 {
		for i, card := range hand {
			rect := rl.NewRectangle(
				offsetX+float32(int32(i)*CARD_WIDTH),
				float32(CARD_POS_Y),
				float32(CARD_WIDTH),
				float32(CARD_HEIGHT),
			)
			window.playerCards = append(window.playerCards, *NewCardModel(card, rect))
		}
	}
}

func (window *Window) updateLastDiscardedCard(card *dm.Card) {
	rectCenter := rl.NewRectangle(
		float32(WINDOW_WIDTH-CARD_WIDTH)/2,
		float32(WINDOW_HEIGHT-CARD_HEIGHT)/2,
		float32(CARD_WIDTH),
		float32(CARD_HEIGHT),
	)
	window.lastDiscardedCard = NewCardModel(card, rectCenter)
}
