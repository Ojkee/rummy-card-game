package window

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	"rummy-card-game/src/game_logic/game_manager"
)

type Window struct {
	mu       sync.Mutex
	clientId int

	onReadyCallback    func(bool)
	sendActionCallback func(cm.ActionMessage)

	running     bool
	stopChannel chan struct{}
	readyButton FuncButton
	isReady     bool
	gameState   game_manager.GAME_STATE

	lastTurnId        int
	playerCards       []CardModel
	discardPile       *dm.CardQueue
	lastDiscardedCard *CardModel
	drawPile          *DrawPileButton
}

func NewWindow() *Window {
	return &Window{
		mu: sync.Mutex{},

		running:     true,
		stopChannel: make(chan struct{}),
		readyButton: *NewFuncButton(),
		isReady:     false,
		gameState:   game_manager.PRE_START,

		lastTurnId:        -1,
		playerCards:       make([]CardModel, 0),
		discardPile:       nil,
		lastDiscardedCard: NewCardModel(nil, DISCARD_PILE_POS),
		drawPile:          NewDrawPileButton(),
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
	mousePos := rl.GetMousePosition()
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		switch window.gameState {
		case game_manager.PRE_START:
			if window.readyButton.isClicked(&mousePos) {
				window.toggleReady()
			}
		case game_manager.IN_GAME:
			window.inRoundManager(&mousePos)
		default:
			break
		}
	}
}

func (window *Window) draw() {
	rl.BeginDrawing()

	rl.ClearBackground(COLOR_DARK_GRAY)

	switch window.gameState {
	case game_manager.PRE_START:
		window.drawWaitingPane()
		break
	case game_manager.IN_GAME:
		window.drawInRound()
		break
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

func (window *Window) SetOnReadyCallback(onReady func(bool)) {
	window.onReadyCallback = onReady
}

func (window *Window) SetActionMessageCallback(sendAction func(cm.ActionMessage)) {
	window.sendActionCallback = sendAction
}

func (window *Window) SetId(id int) {
	window.clientId = id
}

func (window *Window) SetGameState(gameState game_manager.GAME_STATE) {
	window.gameState = gameState
}

func (window *Window) CloseListener() <-chan struct{} {
	return window.stopChannel
}

func (window *Window) Stop() {
	window.unloadGraphics()
	window.running = false
	close(window.stopChannel)
}

func (window *Window) UpdateState(sv cm.StateView) {
	window.mu.Lock()
	defer window.mu.Unlock()

	window.updatePlayerHand(sv.PlayerEntity.Hand)
	window.discardPile = sv.DiscardPile
	window.updateLastDiscardedCard(window.discardPile.SeekBack())
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
	window.lastDiscardedCard = NewCardModel(card, DISCARD_PILE_POS)
}
