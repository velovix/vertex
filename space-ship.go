package main

import (
	"github.com/go-gl/gl/v2.1/gl"
)

var (
	spaceShipColSize           = 80.0
	spaceShipMaxHP             = 1
	spaceShipShootCooldown     = 5.0
	spaceShipBulletSpeed       = 20.0
	spaceShipMass              = 5.0
	spaceShipSpeed             = 5.0
	spaceShipSpawnScalingSpeed = 0.1

	spaceShipInitialInvulnerability = 30

	spaceShipGunTip1 = vertex{30, 25, 0}
	spaceShipGunTip2 = vertex{30, -25, 0}

	spaceShipScale = vertex{50, 50, 50}
	spaceShipModel polyModel
)

// spaceShip is a player object. It is capable of moving around and shooting.
type spaceShip struct {
	uidGenerator

	loc     vertex
	rot     float64
	moveDir vertex

	hp int

	barrel           bool
	shootingCooldown float64
	shooting         bool

	spawnScaling vertex

	invulnerability int

	vulnerability collision
}

func newSpaceShip() *spaceShip {
	ss := new(spaceShip)

	ss.hp = spaceShipMaxHP

	// Initailize collisions
	ss.vulnerability = collision{
		alliance: friendly,
		typ:      vulnerable}
	ss.updateCols()

	ss.invulnerability = spaceShipInitialInvulnerability

	return ss
}

// respawn respawns the player at origin with full health.
func (ss *spaceShip) respawn() {
	ss.hp = spaceShipMaxHP
	ss.loc = vertex{0, 0, 0}
	ss.invulnerability = spaceShipInitialInvulnerability
	ss.spawnScaling = vertex{0, 0, 0}
}

func (ss *spaceShip) tick() []entity {
	var newEnts []entity

	ss.loc.x += (spaceShipSpeed * mainWindow.delta) * ss.moveDir.x
	ss.loc.y += (spaceShipSpeed * mainWindow.delta) * ss.moveDir.y

	if ss.shootingCooldown > 0 {
		ss.shootingCooldown -= 1.0 * mainWindow.delta
	} else if ss.shooting {
		// Shoot a bullet
		newEnts = append(newEnts, ss.shoot())
	}

	if ss.invulnerability > 0 {
		ss.invulnerability--
	}

	// Make the space ship scale up to full size on creation
	if ss.spawnScaling.x < 1.0 {
		ss.spawnScaling.x += spaceShipSpawnScalingSpeed * mainWindow.delta
		ss.spawnScaling.y += spaceShipSpawnScalingSpeed * mainWindow.delta
		ss.spawnScaling.z += spaceShipSpawnScalingSpeed * mainWindow.delta
		if ss.spawnScaling.x > 1.0 {
			ss.spawnScaling = vertex{1, 1, 1}
		}
	}

	ss.updateCols()

	return newEnts
}

func (ss *spaceShip) draw() {
	gl.PushMatrix()

	gl.Translated(ss.loc.x, ss.loc.y, 0.0)
	gl.Rotated(ss.rot, 0, 0, 1)
	gl.Scaled(spaceShipScale.x, spaceShipScale.y, spaceShipScale.z)
	gl.Scaled(ss.spawnScaling.x, ss.spawnScaling.y, ss.spawnScaling.z)

	spaceShipModel.draw()

	gl.PopMatrix()
}

// shoot returns a bullet that is set up to be shot.
func (ss *spaceShip) shoot() *bullet {
	// Find the position of the gun tip where the bullet will be shot
	var rotatedGunTip vertex
	if ss.barrel {
		rotatedGunTip = rotateVertex(vertex{0.0, 0.0, 0.0}, spaceShipGunTip1, ss.rot)
	} else {
		rotatedGunTip = rotateVertex(vertex{0.0, 0.0, 0.0}, spaceShipGunTip2, ss.rot)
	}
	// Switch barrels for next shot
	ss.barrel = !ss.barrel
	ss.shootingCooldown = spaceShipShootCooldown

	return newBullet(bulletPosData{
		loc:     vertex{ss.loc.x + rotatedGunTip.x, ss.loc.y + rotatedGunTip.y, 0.0},
		rot:     ss.rot,
		moveDir: unitVectorFromAngle(ss.rot),
		speed:   spaceShipBulletSpeed,
	}, friendly)
}

func (ss *spaceShip) location() vertex {
	return ss.loc
}

func (ss *spaceShip) collisions() []collision {
	return []collision{
		ss.vulnerability,
	}
}

func (ss *spaceShip) collision(yours, other collision) {
	switch other.typ {
	case harmful:
		if yours.uid() == ss.vulnerability.uid() && other.alliance != friendly {
			if ss.invulnerability <= 0 {
				ss.hp--
			}
		}
	}
}

func (ss *spaceShip) mass() float64 {
	return spaceShipMass
}

func (ss *spaceShip) updateCols() {
	ss.vulnerability.bounding = bounding{
		vertex{
			ss.loc.x - spaceShipColSize/2.0, ss.loc.y - spaceShipColSize/2.0, 0.0},
		vertex{
			ss.loc.x + spaceShipColSize/2.0, ss.loc.y + spaceShipColSize/2.0, 0.0}}
}

func (ss *spaceShip) deletable() bool {
	return ss.hp <= 0
}

func (ss *spaceShip) inFrame() {
}
