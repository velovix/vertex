package main

import "github.com/go-gl/glfw/v3.0/glfw"

// supportedJoysticks is a list of all potential joysticks that can be
// connected.
var supportedJoysticks []glfw.Joystick = []glfw.Joystick{
	glfw.Joystick1, glfw.Joystick2, glfw.Joystick3, glfw.Joystick4}

// connectedJoysticks returns a list of all connected joysticks.
func connectedJoysticks() []glfw.Joystick {
	var output []glfw.Joystick

	for _, joy := range supportedJoysticks {
		if glfw.JoystickPresent(joy) {
			output = append(output, joy)
		}
	}

	return output
}
