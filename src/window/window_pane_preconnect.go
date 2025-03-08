package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (window *Window) preConnectManagetKeyboardInput() {
	if isPasteShortcutClicked() {
		window.pasteFromClipboard(rl.GetClipboardText())
		return
	}
	keyPressed := rl.GetKeyPressed()
	if keyPressed == rl.KeyBackspace && len(window.enteredIp.content) > 0 {
		newContent := window.enteredIp.content[:len(window.enteredIp.content)-1]
		window.enteredIp.UpdateContent(newContent)
		return
	} else if keyPressed == rl.KeyEnter {
		window.connectCallback(window.enteredIp.content)
		return
	} else if len(window.enteredIp.content) >= window.maxIpLen {
		return
	}
	charPressed := rl.GetCharPressed()
	if isValidCharASCII(charPressed) {
		newContent := window.enteredIp.content + string(charPressed)
		window.enteredIp.UpdateContent(newContent)
	}
}

func isPasteShortcutClicked() bool {
	if (rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)) &&
		rl.IsKeyPressed(rl.KeyV) {
		return true
	}
	return false
}

func (window *Window) pasteFromClipboard(text string) {
	var filteredText string
	for _, char := range text {
		if isValidCharASCII(char) {
			filteredText += string(char)
		}
	}
	window.enteredIp.UpdateContent(filteredText)
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
