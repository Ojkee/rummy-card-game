package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

var (
	WINDOW_WIDTH  int32 = 800
	WINDOW_HEIGHT int32 = 800
)

var (
	FONT_SIZE    int32   = 16
	FONT_SPACING float32 = 2
	FONT_PATH    string  = "src/window/assets/fonts/Child-Hood.otf"
	FONT         rl.Font
)

var (
	COLOR_DARK_GRAY    = rl.NewColor(51, 51, 51, 255)
	COLOR_TAUPE        = rl.NewColor(77, 70, 62, 255)
	COLOR_WALNUT_BROWN = rl.NewColor(100, 93, 84, 255)
	COLOR_BEIGE        = rl.NewColor(255, 248, 231, 255)
)

var (
	CARD_WIDTH        int32 = 48
	CARD_HEIGHT       int32 = 96
	CARD_POS_Y        int32 = WINDOW_HEIGHT - CARD_HEIGHT
	CARD_GAP          int32 = 2
	CARD_INNER_WIDTH  int32 = CARD_WIDTH - CARD_GAP*2
	CARD_INNER_HEIGHT int32 = CARD_HEIGHT - CARD_GAP*2
)

var (
	SUIT_WIDTH  int32 = CARD_INNER_WIDTH
	SUIT_HEIGHT int32 = CARD_INNER_WIDTH
)
var RANK_IMGS = make(map[dm.Suit]rl.Texture2D)

var (
	CLUBS_IMG    = rl.LoadImage("src/window/assets/images/Clubs.png")
	DIAMONDS_IMG = rl.LoadImage("src/window/assets/images/Diamonds.png")
	HEARTS_IMG   = rl.LoadImage("src/window/assets/images/Hearts.png")
	SPADES_IMG   = rl.LoadImage("src/window/assets/images/Spades.png")
	JOKER_IMG    = rl.LoadImage("src/window/assets/images/joka.png")
)
