package table_manager

type DEBUG_MODE int

const (
	MELD_HAND_START = iota
)

var DEBUG_MODES = map[DEBUG_MODE]bool{
	MELD_HAND_START: true,
}

const (
	SKIP_MELD_HAND_CARDS = 5
)
