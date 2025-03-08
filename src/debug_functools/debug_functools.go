package debug_functools

import rl "github.com/gen2brain/raylib-go/raylib"

type DEBUG_MODE int

const (
	MELD_HAND_START = iota
	RESET_SERVER
	FAST_QUIT
)

var DEBUG_MODES = map[DEBUG_MODE]bool{
	MELD_HAND_START: false,
	RESET_SERVER:    false,
	FAST_QUIT:       false,
}

const (
	SKIP_MELD_HAND_CARDS = 5
	RESET_SERVER_KEY     = rl.KeyF4
	QUIT_KEY             = rl.KeyQ
)
