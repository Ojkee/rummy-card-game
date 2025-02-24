package window

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
}

func NewAvailableSpot(rect rl.Rectangle, spotType AVAILABLE_SPOT_TYPE) *AvailableSpot {
	return &AvailableSpot{
		rect:     rect,
		spotType: spotType,
	}
}

func (as *AvailableSpot) InRect(mousePos *rl.Vector2) bool {
	return rl.CheckCollisionPointRec(*mousePos, as.rect)
}

func (as *AvailableSpot) Draw() {
	rl.DrawRectangleRec(as.rect, COLOR_HIGHLIGHT_SPOT)
}
