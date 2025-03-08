package window

import rl "github.com/gen2brain/raylib-go/raylib"

func (window *Window) waitingPaneManageKeyboardInput() {
	enterPressed := window.keyboardInputStaticButton(&window.enterNickname, window.maxNicknameLen)
	if window.nickname != window.enterNickname.content {
		window.lockNickname()
	}
	if enterPressed && len(window.enterNickname.content) > 0 {
		window.toggleReady()
		return
	} else if enterPressed && len(window.enterNickname.content) == 0 {
		window.PlaceText("Enter nickname")
	}
}

func (window *Window) lockNickname() {
	window.nickname = window.enterNickname.content
}

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
	window.drawStaticButton(&window.enterNickname)
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
		FONT_SIZE,
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
		FONT_SIZE,
		FONT_SPACING,
		COLOR_BEIGE,
	)
}
