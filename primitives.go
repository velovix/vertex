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

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func (a vertex) distance(b vertex) float64 {
	// Used to use math.Pow here, but this is much faster
	return math.Sqrt(((b.x - a.x) * (b.x - a.x)) + ((b.y - a.y) * (b.y - a.y)))
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

func rotateVertex(center, rotatable vertex, angle float64) vertex {
	s := math.Sin(toRadians(angle))
	c := math.Cos(toRadians(angle))

	rotatable.x -= center.x
	rotatable.y -= center.y

	rotatedX := rotatable.x*c - rotatable.y*s
	rotatedY := rotatable.x*s + rotatable.y*c

	rotatable.x = rotatedX + center.x
	rotatable.y = rotatedY + center.y

	return rotatable
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

func angleFromUnitVector(v vertex) float64 {
	angle := (180.0 / math.Pi) * math.Atan2(v.y, v.x)
	if angle < 0.0 {
		angle = 360 + angle
	}

	return angle
}
