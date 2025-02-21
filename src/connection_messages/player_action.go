package connection_messages

type ACTION_TYPE int

const (
	DRAW_CARD ACTION_TYPE = iota
	DISCARD_CARD
	INITIAL_MELD
	UNSUPPORTED
)

type ActionMessage interface {
	JsonMessage
	GetActionType() ACTION_TYPE
}
