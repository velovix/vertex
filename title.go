package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	// titleRateOfMovement is the rate at which the title should approach its
	// target position.
	titleRateOfMovement = 0.0625

	titleModel polyModel
)

type titleStatus int

const (
	movingTitleStatus titleStatus = iota
	waitingTitleStatus
	transitioningTitleStatus
)

type title struct {
	uidGenerator

	pos, rot             vertex
	targetPos, targetRot vertex

	status titleStatus
}

func newTitle() *title {
	t := new(title)

	t.pos = vertex{0, 200, 30}
	t.rot = vertex{90, 0, 0}
	t.targetPos = vertex{0, 0, -30}
	t.targetRot = vertex{90, 0, 0}

	t.status = movingTitleStatus

	var whatever inputEntity = t
	_ = whatever

	return t
}

func (t *title) draw() {
	gl.PushMatrix()

	gl.Translated(t.pos.x, t.pos.y, t.pos.z)
	gl.Rotated(t.rot.x, 1, 0, 0)
	gl.Rotated(t.rot.y, 0, 1, 0)
	gl.Rotated(t.rot.z, 0, 0, 1)
	gl.Scaled(pulseScale.x, pulseScale.y, pulseScale.z)

	titleModel.draw()

	gl.PopMatrix()
}

func (t *title) tick() []entity {
	t.pos.y += ((t.targetPos.y - t.pos.y) * titleRateOfMovement) * mainWindow.delta
	t.pos.z += ((t.targetPos.z - t.pos.z) * titleRateOfMovement) * mainWindow.delta

	if t.status == transitioningTitleStatus && t.pos.z >= 1 {
		currentReg = &gameplayReg
	}

	return []entity{}
}

func (t *title) joystickAxis(joystick glfw.Joystick, axes []float32) {
}

func (t *title) joystickButton(joystick glfw.Joystick, buttons []byte) {
	for _, btn := range buttons {
		if btn == 1 {
			stopMusic()
			resetPulse()
			playMusic(gameplayMusic, 0)
			t.targetPos = vertex{0, 0, 5}
			t.status = transitioningTitleStatus
		}
	}
}

func (t *title) location() vertex {
	return t.pos
}

func (t *title) collisions() []collision {
	return []collision{}
}

func (t *title) collision(yours, other collision) {
}

func (t *title) mass() float64 {
	return 0
}

func (t *title) deletable() bool {
	return false
}
