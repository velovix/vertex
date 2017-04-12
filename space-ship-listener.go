package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	joystickSensitivity = 0.1
)

type spaceShipListener struct {
	*spaceShip

	joystick glfw.Joystick
}

func newSpaceShipListener(joystick glfw.Joystick) *spaceShipListener {
	ssl := new(spaceShipListener)
	ssl.spaceShip = newSpaceShip()
	ssl.joystick = joystick
	return ssl
}

func (ssl *spaceShipListener) joystickAxis(joystick glfw.Joystick, axes []float32) {
	if ssl.joystick != joystick {
		return
	}
	if len(axes) < 4 {
		// TODO(velovix): Do something to make sure this never happens
		panic("not enough axes on joystick " + string(joystick))
	}

	// Use controller axes 0 and 1 for movement
	if math.Abs(float64(axes[0])) > joystickSensitivity ||
		math.Abs(float64(axes[1])) > joystickSensitivity {
		ssl.moveDir.x = float64(axes[0])
		ssl.moveDir.y = float64(-axes[1])
	} else {
		ssl.moveDir.x = 0
		ssl.moveDir.y = 0
	}

	// Use controller axes 2 and 3 for aiming
	if math.Abs(float64(axes[3])) > joystickSensitivity ||
		math.Abs(float64(axes[4])) > joystickSensitivity {
		ssl.rot = angleFromUnitVector(vertex{
			float64(axes[3]),
			float64(-axes[4]),
			0.0})
	}
}

func (ssl *spaceShipListener) joystickButton(joystick glfw.Joystick, buttons []byte) {
	if ssl.joystick != joystick {
		return
	}
	if len(buttons) < 6 {
		// TODO(velovix): Do something ot make sure this never happens
		panic("not enough buttons on joystick " + string(joystick))
	}

	if buttons[5] == 1 {
		ssl.shooting = true
	} else if buttons[5] == 0 {
		ssl.shooting = false
	}
}
