package main

import (
	"github.com/go-gl/gl/v2.1/gl"
)

var (
	spaceShipColSize       = 80.0
	spaceShipMaxHP         = 1
	spaceShipShootCooldown = 5.0
	spaceShipMass          = 5.0
	spaceShipSpeed         = 5.0

	spaceShipGunTip1 = vertex{30, 25, 0}
	spaceShipGunTip2 = vertex{30, -25, 0}

	spaceShipScale = vertex{50, 50, 50}
	spaceShipModel polyModel
)

// spaceShip is a player object. It is capable of moving around and shooting.
type spaceShip struct {
	uidGenerator

	x, y, rot float64
	moveDir   vertex

	hp int

	barrel           int
	shootingCooldown float64
	shooting         bool

	vulnerability collision
}

func newSpaceShip() *spaceShip {
	ss := new(spaceShip)

	// Initailize collisions
	ss.vulnerability = collision{
		alliance: friendly,
		typ:      vulnerable}
	ss.updateCols()

	return ss
}

func (ss *spaceShip) tick() []entity {
	return []entity{}
}

func (ss *spaceShip) draw() {
	gl.PushMatrix()

	gl.Translated(ss.x, ss.y, 0.0)
	gl.Rotated(ss.rot, 0, 0, 1)
	gl.Scaled(spaceShipScale.x, spaceShipScale.y, spaceShipScale.z)

	spaceShipModel.draw()

	gl.PopMatrix()
}

func (ss *spaceShip) location() vertex {
	return vertex{x: ss.x, y: ss.y}
}

func (ss *spaceShip) collisions() []collision {
	return []collision{
		ss.vulnerability,
	}
}

func (ss *spaceShip) collision(yours, other collision) {
}

func (ss *spaceShip) mass() float64 {
	return spaceShipMass
}

func (ss *spaceShip) updateCols() {
	ss.vulnerability.bounding = bounding{
		vertex{
			ss.x - spaceShipColSize/2.0, ss.y - spaceShipColSize/2.0, 0.0},
		vertex{
			ss.x + spaceShipColSize/2.0, ss.y + spaceShipColSize/2.0, 0.0}}
}

func (ss *spaceShip) deletable() bool {
	return false
}
