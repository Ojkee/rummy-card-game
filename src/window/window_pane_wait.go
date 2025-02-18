package window

import rl "github.com/gen2brain/raylib-go/raylib"

func (window *Window) toggleReady() {
	window.isReady = !window.isReady
	window.onReadyCallback(window.isReady)
}

func (window *Window) drawWaitingPane() {
	if window.isReady {
		window.readyButton.color = COLOR_BUTTON_READY
		window.readyButton.content = "Ready"
	} else {
		window.readyButton.color = COLOR_BUTTON_NOT_READY
		window.readyButton.content = "Not ready"
	}
	window.drawReadyButton()
	window.drawReadyState()
	window.drawReadyInfo()
}

func (window *Window) drawReadyButton() {
	var roundness float32 = 0.75
	var segments int32 = 10
	rl.DrawRectangleRounded(
		window.readyButton.rect,
		roundness,
		segments,
		window.readyButton.color,
	)
}

func (window *Window) drawReadyState() {
	contentVec := GetTextVec(window.readyButton.content)
	rl.DrawTextEx(
		FONT,
		window.readyButton.content,
		rl.NewVector2(
			window.readyButton.rect.X+(window.readyButton.rect.Width-contentVec.X)/2,
			window.readyButton.rect.Y+(window.readyButton.rect.Height-contentVec.Y)/2,
		),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
}

func (window *Window) drawReadyInfo() {
	info := "You are: "
	infoVec := GetTextVec(info)
	rl.DrawTextEx(
		FONT,
		info,
		rl.NewVector2(
			window.readyButton.rect.X+(window.readyButton.rect.Width-infoVec.X)/2,
			window.readyButton.rect.Y-infoVec.Y,
		),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
}
