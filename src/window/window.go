package window

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	cm "rummy-card-game/src/connection_messages"
	dm "rummy-card-game/src/game_logic/deck_manager"
	gm "rummy-card-game/src/game_logic/game_manager"
)

type Window struct {
	mu       sync.Mutex
	clientId int

	onReadyCallback    func(bool)
	sendActionCallback func(cm.ActionMessage)

	readyButton       FuncButton
	discardButton     FuncButton
	lockSetButton     FuncButton
	initialMeldButton FuncButton

	running     bool
	stopChannel chan struct{}
	isReady     bool
	gameState   gm.GAME_STATE

	isDragging         bool
	startDragPos       rl.Vector2
	currentDragCardIdx int
	dragCardStartRec   rl.Rectangle

	currentTurnId     int
	playerCards       []CardModel
	discardPile       *dm.CardQueue
	lastDiscardedCard *CardModel
	drawPile          *DrawPileButton

	displayText string
	displayTime float32

	lockedSequencesIds map[int]bool

	tableSequences []SequenceModel
	availableSpots []AvailableSpot
}

func NewWindow() *Window {
	_lockedSequencesIds := make(map[int]bool)
	for i := range MAX_LOCKS {
		_lockedSequencesIds[i] = false
	}
	return &Window{
		mu: sync.Mutex{},

		readyButton: *NewFuncButton(
			rl.NewRectangle(
				float32(WINDOW_WIDTH-READY_BUTTON_WIDTH)/2,
				float32(WINDOW_HEIGHT-READY_BUTTON_HEIGHT)/2,
				READY_BUTTON_WIDTH,
				READY_BUTTON_HEIGHT,
			),
			"Not ready",
		),
		discardButton: *NewFuncButton(
			rl.NewRectangle(0, 0, DISCARD_BUTTON_WIDTH, DISCARD_BUTTON_HEIGHT),
			"Discard",
		),
		lockSetButton: *NewFuncButton(
			rl.NewRectangle(
				float32(WINDOW_WIDTH-LOCK_SEQUENCE_BUTTON_WIDTH)/2,
				float32(WINDOW_HEIGHT-CARD_HEIGHT)-LOCK_SEQUENCE_BUTTON_HEIGHT-24,
				LOCK_SEQUENCE_BUTTON_WIDTH,
				LOCK_SEQUENCE_BUTTON_HEIGHT,
			),
			"Lock",
		),
		initialMeldButton: *NewFuncButton(
			rl.NewRectangle(
				float32(WINDOW_WIDTH-LOCK_SEQUENCE_BUTTON_WIDTH)/2,
				float32(WINDOW_HEIGHT-CARD_HEIGHT)-2*LOCK_SEQUENCE_BUTTON_HEIGHT-28,
				LOCK_SEQUENCE_BUTTON_WIDTH,
				LOCK_SEQUENCE_BUTTON_HEIGHT,
			),
			"Meld",
		),

		running:     true,
		stopChannel: make(chan struct{}),
		isReady:     false,
		gameState:   gm.PRE_START,

		isDragging:         false,
		currentDragCardIdx: -1,

		currentTurnId:     -1,
		playerCards:       make([]CardModel, 0),
		discardPile:       nil,
		lastDiscardedCard: NewCardModel(nil, DISCARD_PILE_POS, false),
		drawPile:          NewDrawPileButton(),

		lockedSequencesIds: _lockedSequencesIds,
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
		window.startDragPos = mousePos
		window.isDragging = false
	}
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		if !window.isDragging && rl.Vector2Distance(window.startDragPos, mousePos) > 5 {
			window.isDragging = true
		}
		if window.isDragging {
			window.handleMouseDrag(&mousePos)
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if !window.isDragging {
			window.handleMouseClicked(&mousePos)
		}
		window.isDragging = false
		window.rearrangeNewCardPosX()
		window.currentDragCardIdx = -1
	}
}

func (window *Window) draw() {
	rl.BeginDrawing()
	rl.ClearBackground(COLOR_DARK_GRAY)

	switch window.gameState {
	case gm.PRE_START:
		window.drawWaitingPane()
		break
	case gm.IN_GAME:
		window.drawInRound()
		break
	}

	rl.EndDrawing()
}

func (window *Window) initGraphics() {
	window.resizeImages()
	window.loadImagesMap()
	FONT = rl.LoadFontEx(FONT_PATH, FONT_SIZE, nil, 96)
}

func (window *Window) resizeImages() {
	rl.ImageResize(CLUBS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(DIAMONDS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(HEARTS_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(SPADES_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(JOKER_IMG, SUIT_WIDTH, SUIT_HEIGHT)
	rl.ImageResize(CLUBS_IMG_SMALL, SUIT_WIDTH_SMALL, SUIT_HEIGHT_SMALL)
	rl.ImageResize(DIAMONDS_IMG_SMALL, SUIT_WIDTH_SMALL, SUIT_HEIGHT_SMALL)
	rl.ImageResize(HEARTS_IMG_SMALL, SUIT_WIDTH_SMALL, SUIT_HEIGHT_SMALL)
	rl.ImageResize(SPADES_IMG_SMALL, SUIT_WIDTH_SMALL, SUIT_HEIGHT_SMALL)
	rl.ImageResize(JOKER_IMG_SMALL, SUIT_WIDTH_SMALL, SUIT_HEIGHT_SMALL)
}

func (window *Window) loadImagesMap() {
	RANK_IMGS[dm.CLUBS] = rl.LoadTextureFromImage(CLUBS_IMG)
	RANK_IMGS[dm.DIAMONDS] = rl.LoadTextureFromImage(DIAMONDS_IMG)
	RANK_IMGS[dm.HEARTS] = rl.LoadTextureFromImage(HEARTS_IMG)
	RANK_IMGS[dm.SPADES] = rl.LoadTextureFromImage(SPADES_IMG)
	RANK_IMGS[dm.ANY] = rl.LoadTextureFromImage(JOKER_IMG)
	RANK_IMGS_SMALL[dm.CLUBS] = rl.LoadTextureFromImage(CLUBS_IMG_SMALL)
	RANK_IMGS_SMALL[dm.DIAMONDS] = rl.LoadTextureFromImage(DIAMONDS_IMG_SMALL)
	RANK_IMGS_SMALL[dm.HEARTS] = rl.LoadTextureFromImage(HEARTS_IMG_SMALL)
	RANK_IMGS_SMALL[dm.SPADES] = rl.LoadTextureFromImage(SPADES_IMG_SMALL)
	RANK_IMGS_SMALL[dm.ANY] = rl.LoadTextureFromImage(JOKER_IMG_SMALL)
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

func (window *Window) SetClientId(id int) {
	window.clientId = id
}

func (window *Window) SetGameState(gameState gm.GAME_STATE) {
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
	window.currentTurnId = sv.TurnPlayerId
	window.updateTableSequences(sv.TableSequences)
}

func (window *Window) PlaceText(text string) {
	window.displayText = text
	window.displayTime = TIME_ON_SCREEN
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
			window.playerCards = append(window.playerCards, *NewCardModel(card, rect, false))
		}
	}
}

func (window *Window) updateLastDiscardedCard(card *dm.Card) {
	window.lastDiscardedCard = NewCardModel(card, DISCARD_PILE_POS, false)
}

func (window *Window) updateTableSequences(sequences []gm.Sequence) {
	tableSequencesNew := make([]SequenceModel, 0)
	for i, sequence := range sequences {
		firstCardPos := rl.NewVector2(
			TABLE_X,
			float32(i)*SEQUENCE_CARD_HEIGHT+TABLE_Y,
		)
		sequenceModel := NewSequenceModel(
			&sequence,
			firstCardPos,
		)
		tableSequencesNew = append(tableSequencesNew, *sequenceModel)
	}
	window.tableSequences = tableSequencesNew
}

func (window *Window) drawStaticButton(fbutton *FuncButton) {
	rectInner := fbutton.rect
	rectInner.X += 2
	rectInner.Y += 2
	rectInner.Width -= 4
	rectInner.Height -= 4
	rl.DrawRectangleRounded(fbutton.rect, 0.5, 10, COLOR_BEIGE)
	rl.DrawRectangleRounded(rectInner, 0.5, 10, COLOR_DARK_GRAY)

	contentSize := GetTextVec(fbutton.content)
	rl.DrawTextEx(
		FONT,
		fbutton.content,
		rl.NewVector2(
			fbutton.rect.X+(fbutton.rect.Width-contentSize.X)/2,
			fbutton.rect.Y+(fbutton.rect.Height-contentSize.Y)/2,
		),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
}

func (window *Window) handleMouseClicked(mousePos *rl.Vector2) {
	switch window.gameState {
	case gm.PRE_START:
		if window.readyButton.InRect(mousePos) {
			window.toggleReady()
		}
	case gm.IN_GAME:
		window.inRoundManagerClick(mousePos)
	default:
		break
	}
}

func (window *Window) handleMouseDrag(mousePos *rl.Vector2) {
	switch window.gameState {
	case gm.IN_GAME:
		window.inRoundManagerDrag(mousePos)
	default:
		break
	}
}
