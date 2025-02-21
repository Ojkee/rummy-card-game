package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

const (
	WINDOW_WIDTH  int32 = 800
	WINDOW_HEIGHT int32 = 800
)

var (
	FONT_SIZE    int32   = 24
	FONT_SPACING float32 = 2
	FONT_PATH    string  = "src/window/assets/fonts/Child-Hood.otf"
	FONT         rl.Font
)

var TIME_ON_SCREEN float32 = 3.0

var (
	COLOR_DARK_GRAY        = rl.NewColor(51, 51, 51, 255)
	COLOR_TAUPE            = rl.NewColor(77, 70, 62, 255)
	COLOR_WALNUT_BROWN     = rl.NewColor(100, 93, 84, 255)
	COLOR_BEIGE            = rl.NewColor(255, 248, 231, 255)
	COLOR_BUTTON_READY     = rl.NewColor(81, 123, 74, 255)
	COLOR_BUTTON_NOT_READY = rl.NewColor(144, 55, 55, 255)
	COLOR_LOCK_COLOR_1     = rl.NewColor(255, 0, 0, 255)
	COLOR_LOCK_COLOR_2     = rl.NewColor(0, 255, 0, 255)
	COLOR_LOCK_COLOR_3     = rl.NewColor(0, 0, 255, 255)
	COLOR_LOCK_COLOR_4     = rl.NewColor(255, 0, 255, 255)
)

var LOCK_COLORS = map[int]rl.Color{
	0: COLOR_LOCK_COLOR_1,
	1: COLOR_LOCK_COLOR_2,
	2: COLOR_LOCK_COLOR_3,
	3: COLOR_LOCK_COLOR_4,
}

var MAX_LOCKS = len(LOCK_COLORS)

// BUTTONS
const (
	READY_BUTTON_WIDTH          = 128
	READY_BUTTON_HEIGHT         = 64
	IN_GAME_BUTTON_WIDHT        = 128
	IN_GAME_BUTTON_HEIGHT       = 64
	DISCARD_BUTTON_WIDTH        = 64
	DISCARD_BUTTON_HEIGHT       = 32
	LOCK_SEQUENCE_BUTTON_WIDTH  = 64
	LOCK_SEQUENCE_BUTTON_HEIGHT = 32
)

// CARDS
const (
	CARD_WIDTH        int32 = 48
	CARD_HEIGHT       int32 = 96
	CARD_POS_Y        int32 = WINDOW_HEIGHT - CARD_HEIGHT
	CARD_GAP          int32 = 2
	CARD_INNER_WIDTH  int32 = CARD_WIDTH - CARD_GAP*2
	CARD_INNER_HEIGHT int32 = CARD_HEIGHT - CARD_GAP*2
)

const (
	PILES_OFFSET float32 = float32(CARD_WIDTH)/2 + 10
)

var (
	DISCARD_PILE_POS = rl.NewRectangle(
		float32(WINDOW_WIDTH-CARD_WIDTH)/2-PILES_OFFSET,
		float32(WINDOW_HEIGHT-CARD_HEIGHT)/2,
		float32(CARD_WIDTH),
		float32(CARD_HEIGHT),
	)
	DRAW_PILE_POS = rl.NewRectangle(
		float32(WINDOW_WIDTH-CARD_WIDTH)/2+PILES_OFFSET,
		float32(WINDOW_HEIGHT-CARD_HEIGHT)/2,
		float32(CARD_WIDTH),
		float32(CARD_HEIGHT),
	)
)

const (
	SUIT_WIDTH  int32 = CARD_INNER_WIDTH
	SUIT_HEIGHT int32 = CARD_INNER_WIDTH
)

var RANK_IMGS = make(map[dm.Suit]rl.Texture2D)

var (
	CLUBS_IMG    = rl.LoadImage("src/window/assets/images/Clubs.png")
	DIAMONDS_IMG = rl.LoadImage("src/window/assets/images/Diamonds.png")
	HEARTS_IMG   = rl.LoadImage("src/window/assets/images/Hearts.png")
	SPADES_IMG   = rl.LoadImage("src/window/assets/images/Spades.png")
	JOKER_IMG    = rl.LoadImage("src/window/assets/images/Joker.png")
)
