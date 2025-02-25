package game_manager

import rl "github.com/gen2brain/raylib-go/raylib"

type AVAILABLE_SPOT_TYPE int

const (
	REPLACE_JOKER AVAILABLE_SPOT_TYPE = iota
	ADD_BEGIN
	ADD_END
)

type AvailableSpot struct {
	rect     rl.Rectangle
	spotType AVAILABLE_SPOT_TYPE
	color    rl.Color
	sequence Sequence
}

func NewAvailableSpot(
	rect rl.Rectangle,
	spotType AVAILABLE_SPOT_TYPE,
	color rl.Color,
	sequence Sequence,
) *AvailableSpot {
	return &AvailableSpot{
		rect:     rect,
		spotType: spotType,
		color:    color,
		sequence: sequence,
	}
}

func (as *AvailableSpot) InRect(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, as.rect)
}

func (as *AvailableSpot) Draw() {
	rl.DrawRectangleRec(as.rect, as.color)
}

func (as *AvailableSpot) GetSpotType() AVAILABLE_SPOT_TYPE {
	return as.spotType
}

func (as *AvailableSpot) GetSequence() Sequence {
	return as.sequence
}
