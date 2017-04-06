package main

import "github.com/go-gl/glfw/v3.0/glfw"

const (
	defaultWidth  = 1280
	defaultHeight = 970
)

type window struct {
	width  int
	height int
	delta  float64

	lastFrame float64
	glfw      *glfw.Window
}

var mainWindow window = window{
	width:  defaultWidth,
	height: defaultHeight,
	delta:  1.0,
}
