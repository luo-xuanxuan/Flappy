package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// Call runtime.LockOSThread() in the main function, too, if GLFW event handling
	// or drawing is done in separate goroutines.
	runtime.LockOSThread()
}

func handleGamepadInput() {
	for jid := glfw.Joystick1; jid <= glfw.JoystickLast; jid++ {
		if glfw.Joystick(jid).Present() {
			buttons := glfw.Joystick(jid).GetButtons()
			processFlappy(buttons)
		}
	}
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(fmt.Sprintf("Failed to initialize GLFW: %v", err))
	}
	defer glfw.Terminate()

	// GLFW window creation is not necessary if only polling for gamepad
	window, err := glfw.CreateWindow(640, 480, "GLFW Gamepad Example", nil, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create window: %v", err))
	}

	window.MakeContextCurrent()
	glfw.PollEvents() // Poll for events

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	// Main loop
	for !window.ShouldClose() {
		handleGamepadInput()
		glfw.PollEvents() // Poll for events such as gamepad input
	}

}
