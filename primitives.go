package main

type direction int

const (
	left direction = iota
	right
	up
	down
)

type vertex struct {
	x, y, z float64
}

type bounding struct {
	a, b vertex
}

func unitVectorFromAngle(angle float64) vertex {
	v := vertex{}

	if angle <= 90 {
		v.x = 1.0 - (angle / 90.0)
		v.y = 1.0 - v.x
	} else if angle > 90 && angle <= 180 {
		v.x = -((angle - 90.0) / 90.0)
		v.y = 1.0 + v.x
	} else if angle > 180.0 && angle <= 270.0 {
		v.x = -1.0 + ((angle - 180.0) / 90.0)
		v.y = -1.0 - v.x
	} else if angle > 270 {
		v.x = ((angle - 270.0) / 90.0)
		v.y = -1.0 + v.x
	}

	return v
}
