package game_manager

type GAME_STATE int

const (
	PRE_START GAME_STATE = iota
	IN_GAME
	FINISHED
)
