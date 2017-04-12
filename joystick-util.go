package main

import (
	"bytes"

	"github.com/go-gl/glfw/v3.2/glfw"
)

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

// lastButtons contains the button values of all joysticks last check.
var lastButtons map[glfw.Joystick][]byte

// pollJoysticks polls for Joystick events ant notifies the current entity
// registry.
func pollJoysticks() error {
	for _, joy := range connectedJoysticks() {
		axes, err := glfw.GetJoystickAxes(joy)
		if err != nil {
			return err
		}
		currentReg.joystickAxisEvent(joy, axes)

		buttons, err := glfw.GetJoystickButtons(joy)
		if err != nil {
			return err
		}
		// Only notify about buttons if they've changed
		if !bytes.Equal(buttons, lastButtons[joy]) {
			currentReg.joystickButtonEvent(joy, buttons)
			lastButtons[joy] = buttons
		}
	}

	return nil
}

func init() {
	lastButtons = make(map[glfw.Joystick][]byte)
}
