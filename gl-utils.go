package main

import "github.com/go-gl/gl/v2.1/gl"

// rotateOnPoint executes OpenGL commands to rotate around the given x and y
// point.
func rotateOnPoint(rot float64, loc vertex) {
	gl.Translated(loc.x, loc.y, loc.z)
	gl.Rotated(rot, 0.0, 0.0, 1.0)
	gl.Translated(-loc.x, -loc.y, -loc.z)
}

// fPtr converts the given []float32 to a *float32-casted unsafe pointer. This
// is for OpenGL functions that need data in ths format.
func fPtr(data []float32) *float32 {
	return (*float32)(gl.Ptr(data))
}
