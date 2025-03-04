package window

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type timeGraphics struct {
	startDuration float32
	duration      float32
	color         rl.Color
}

func (tg *timeGraphics) GetDuration() float32 {
	return tg.duration
}

func (tg *timeGraphics) DecrementDuraton(time float32) {
	tg.duration -= time
}

func (tg *timeGraphics) UpdateOpacity() {
	if tg.duration <= 0 {
		return
	}
	f := func(x float64) float64 {
		return math.Pow(x, 3)
	}
	ratio := f(float64(tg.duration / tg.startDuration))
	opacity := uint8(ratio * 255.0)
	tg.color.A = opacity
}

type TextTimeGraphics struct {
	timeGraphics
	text string
	pos  rl.Vector2
}

func NewTextTimeGraphics(
	text string,
	duration float32,
	pos rl.Vector2,
	color rl.Color,
) *TextTimeGraphics {
	return &TextTimeGraphics{
		timeGraphics: timeGraphics{
			startDuration: duration,
			duration:      duration,
			color:         color,
		},
		text: text,
		pos:  pos,
	}
}

func (ttg *TextTimeGraphics) Draw() {
	rl.DrawTextEx(
		FONT,
		ttg.text,
		ttg.pos,
		float32(FONT_SIZE),
		FONT_SPACING,
		ttg.color,
	)
}

type RectTimeGraphics struct {
	timeGraphics
	seqId int
	rect  rl.Rectangle
}

func NewRectTimeGraphics(
	duration float32,
	seqId int,
	rect rl.Rectangle,
	color rl.Color,
) *RectTimeGraphics {
	return &RectTimeGraphics{
		timeGraphics: timeGraphics{
			startDuration: duration,
			duration:      duration,
			color:         color,
		},
		rect: rect,
	}
}

func (rtg *RectTimeGraphics) Draw() {
	rl.DrawRectangleRec(rtg.rect, rtg.color)
}
