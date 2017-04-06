package main

import "math"

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

func (a vertex) distance(b vertex) float64 {
	return math.Sqrt(math.Pow(b.x-a.x, 2) + math.Pow(b.y-a.y, 2))
}

func (a vertex) subtract(b vertex) vertex {
	return vertex{a.x - b.x, a.y - b.y, a.z - b.z}
}

func normalize(v vertex) vertex {
	sum := math.Abs(v.x) + math.Abs(v.y) + math.Abs(v.z)
	if sum == 0 {
		sum = 1
	}

	return vertex{
		x: v.x / sum,
		y: v.y / sum,
		z: v.z / sum}
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
