package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (window *Window) preConnectManageKeyboardInput() {
	enterPressed := window.keyboardInputStaticButton(&window.enteredIp, window.maxIpLen)
	if enterPressed {
		window.connectCallback(window.enteredIp.content)
		return
	}
}

func isPasteShortcutClicked() bool {
	if (rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)) &&
		rl.IsKeyPressed(rl.KeyV) {
		return true
	}
	return false
}

func (window *Window) pasteFromClipboard(text string, fb *FuncButton) {
	var filteredText string
	for _, char := range text {
		if isValidCharASCII(char) {
			filteredText += string(char)
		}
	}
	fb.UpdateContent(filteredText)
}

func isValidCharASCII(c int32) bool {
	return c >= 32 && c <= 126
}

func (window *Window) preConnectManagerClick(mousePos *rl.Vector2) {
	if window.connectButton.InRect(mousePos) {
		window.connectCallback(window.enteredIp.content)
	}
}

func (window *Window) drawPreConnectPane() {
	window.drawStaticButton(&window.connectButton)
	window.drawStaticButton(&window.enteredIp)
	window.drawIpInfo()
}

func (window *Window) drawIpInfo() {
	info := "Enter ip"
	infoVec := GetTextVec(info)
	pos := rl.NewVector2(
		window.enteredIp.rect.X+(window.enteredIp.rect.Width-infoVec.X)/2,
		window.enteredIp.rect.Y-infoVec.Y,
	)
	rl.DrawTextEx(
		FONT,
		info,
		pos,
		FONT_SIZE,
		FONT_SPACING,
		COLOR_BEIGE,
	)
}
