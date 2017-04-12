package main

import "github.com/go-gl/glfw/v3.2/glfw"

const (
	defaultWidth  = 1280
	defaultHeight = 970

	targetSPF = 1.0 / 60.0
)

type window struct {
	width  int
	height int
	delta  float64

	lastFrame float64
	glfw      *glfw.Window
}

func (w *window) calcDelta() {
	if w.lastFrame == 0.0 {
		// There hasn't been a frame to compare this one to, assume delta of 1
		w.lastFrame = glfw.GetTime()
		w.delta = 1.0
	} else {
		currFrame := glfw.GetTime()

		w.delta = (currFrame - w.lastFrame) / targetSPF

		w.lastFrame = currFrame
	}
}

var mainWindow window = window{
	width:  defaultWidth,
	height: defaultHeight,
	delta:  1.0,
}
