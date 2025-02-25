package connection_messages

type ACTION_TYPE int

const (
	DRAW_CARD ACTION_TYPE = iota
	DISCARD_CARD
	INITIAL_MELD
	REARRANGE_CARDS
	UPDATE_TABLE_SEQUNCE
	UNSUPPORTED
)

type ActionMessage interface {
	JsonMessage
	GetActionType() ACTION_TYPE
}
