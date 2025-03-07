package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (window *Window) preConnectManagetKeyboardInput() {
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
	isValidChar := func(c int32) bool {
		return c >= 32 && c <= 126
	}
	charPressed := rl.GetCharPressed()
	if isValidChar(charPressed) {
		newContent := window.enteredIp.content + string(charPressed)
		window.enteredIp.UpdateContent(newContent)
	}
}

func (window *Window) preConnectManagerClick(mousePos *rl.Vector2) {
	if window.connectButton.InRect(mousePos) {
		window.connectCallback(window.enteredIp.content)
	}
}

func (window *Window) drawPreConnectPane() {
	window.drawStaticButton(&window.connectButton)
	window.drawStaticButton(&window.enteredIp)
}
