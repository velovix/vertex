package main

import "time"

const pulseInterval time.Duration = time.Millisecond * 463

var (
	pulseTick  <-chan time.Time
	pulseScale vertex
)

func resetPulse() {
	pulseTick = time.Tick(pulseInterval)
	pulseScale = vertex{1, 1, 1}
}

func updatePulse() {
	select {
	case <-pulseTick:
		pulseScale = vertex{1.1, 1.1, 1.1}
	default:
	}

	if pulseScale.x > 1.0 {
		pulseScale.x -= 0.01 * mainWindow.delta
		pulseScale.y -= 0.01 * mainWindow.delta
		pulseScale.z -= 0.01 * mainWindow.delta
	}

	if pulseScale.x < 1.0 {
		pulseScale = vertex{1.0, 1.0, 1.0}
	}
}
