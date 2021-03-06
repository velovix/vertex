package main

import (
	"math/rand"

	"github.com/go-gl/gl/v2.1/gl"
)

var (
	fanEnemySize              = 60.0
	fanEnemyMoveBoxSize       = 13.0 // TODO(velovix): Remember what this represents
	fanEnemyMaxHP             = 5
	fanEnemySpeed             = 2.5
	fanEnemyRotSpeed          = 1.0
	fanEnemyMass              = 1.5
	fanEnemySpawnScalingSpeed = 0.05

	fanEnemyScale = vertex{25, 25, 25}
	fanEnemyModel polyModel
)

// fanEnemy is a basic enemy that wanders around the stage and hurts the player
// on impact.
type fanEnemy struct {
	uidGenerator

	loc     vertex
	rot     float64
	moveDir vertex

	spawnScaling vertex

	hp int

	harmful, vulnerability collision
	bouncers               map[direction]*collision
}

func newFanEnemy(loc vertex) *fanEnemy {
	fe := new(fanEnemy)

	fe.loc = loc
	fe.hp = fanEnemyMaxHP

	// Start moving around in a random direction
	fe.moveDir = unitVectorFromAngle(360.0 * rand.Float64())

	// Initialize collisions
	fe.harmful = collision{
		alliance: unfriendly,
		typ:      harmful}
	fe.vulnerability = collision{
		alliance: unfriendly,
		typ:      vulnerable}
	fe.bouncers = make(map[direction]*collision)
	for _, dir := range []direction{up, down, left, right} {
		fe.bouncers[dir] = &collision{
			alliance: unfriendly,
			typ:      bouncer}
	}
	fe.updateCols()

	return fe
}

func (fe *fanEnemy) tick() []entity {
	// Move around
	fe.loc.x += (fanEnemySpeed * mainWindow.delta) * fe.moveDir.x
	fe.loc.y += (fanEnemySpeed * mainWindow.delta) * fe.moveDir.y

	// Do a constant rotation for aesthetics
	fe.rot += fanEnemyRotSpeed * mainWindow.delta
	if fe.rot > 360 {
		fe.rot = 0
	}

	// Make the fan enemy scale up to full size on creation
	if fe.spawnScaling.x < 1.0 {
		fe.spawnScaling.x += fanEnemySpawnScalingSpeed * mainWindow.delta
		fe.spawnScaling.y += fanEnemySpawnScalingSpeed * mainWindow.delta
		fe.spawnScaling.z += fanEnemySpawnScalingSpeed * mainWindow.delta
		if fe.spawnScaling.x > 1.0 {
			fe.spawnScaling = vertex{1, 1, 1}
		}
	}

	fe.updateCols()

	return []entity{}
}

func (fe *fanEnemy) draw() {
	gl.PushMatrix()

	rotateOnPoint(fe.rot, fe.loc)
	gl.Translated(fe.loc.x, fe.loc.y, 0.0)
	gl.Scaled(fanEnemyScale.x, fanEnemyScale.y, fanEnemyScale.z)
	gl.Scaled(pulseScale.x, pulseScale.y, pulseScale.z)
	gl.Scaled(fe.spawnScaling.x, fe.spawnScaling.y, fe.spawnScaling.z)

	fanEnemyModel.draw()

	gl.PopMatrix()
}

func (fe *fanEnemy) collision(yours, other collision) {
	switch other.typ {
	case bouncer:
		// Reverse direction if the entity bounces on something
		if is, dir := fe.isBouncer(yours); is {
			if dir == up || dir == down {
				fe.moveDir.y = -fe.moveDir.y
				if fe.moveDir.y > 0 {
					fe.loc.y += 2.5
				} else {
					fe.loc.y -= 2.5
				}
			} else {
				fe.moveDir.x = -fe.moveDir.x
				if fe.moveDir.x > 0 {
					fe.loc.x += 2.5
				} else {
					fe.loc.x -= 2.5
				}
			}
		}
	case harmful:
		// Destroy fan enemies that are inside each other
		if yours.uid() == fe.harmful.uid() && other.alliance == unfriendly {
			fe.hp = 0
		}
		// Take damange if under opposing fire
		if yours.uid() == fe.vulnerability.uid() && other.alliance != unfriendly {
			fe.hp--
		}
	case deadly:
		if yours.uid() == fe.harmful.uid() && other.alliance != unfriendly {
			fe.hp = 0
		}
	}

}

// isBouncer takes in a collision and returns true and the corresponding
// direction if it is one of the Fan Enemy's bouncer collisions.
func (fe *fanEnemy) isBouncer(col collision) (bool, direction) {
	for key, val := range fe.bouncers {
		if val.uid() == col.uid() {
			return true, key
		}
	}
	return false, up
}

// updateCols updates the location of collisions to the position of the entity.
func (fe *fanEnemy) updateCols() {
	fe.harmful.bounding = bounding{
		vertex{fe.loc.x - fanEnemySize/2, fe.loc.y - fanEnemySize/2, 0.0},
		vertex{fe.loc.x + fanEnemySize/2, fe.loc.y + fanEnemySize/2, 0.0}}
	fe.vulnerability.bounding = fe.harmful.bounding

	fe.bouncers[up].bounding =
		bounding{
			vertex{fe.loc.x - fanEnemySize/2, fe.loc.y + fanEnemySize/2, 0.0},
			vertex{fe.loc.x + fanEnemySize/2, fe.loc.y + fanEnemySize/2 + fanEnemyMoveBoxSize, 0.0}}
	fe.bouncers[down].bounding =
		bounding{
			vertex{fe.loc.x - fanEnemySize/2, fe.loc.y - fanEnemySize/2 - fanEnemyMoveBoxSize, 0.0},
			vertex{fe.loc.x + fanEnemySize/2, fe.loc.y - fanEnemySize/2, 0.0}}
	fe.bouncers[left].bounding =
		bounding{
			vertex{fe.loc.x - fanEnemySize/2 - fanEnemyMoveBoxSize, fe.loc.y - fanEnemySize/2, 0.0},
			vertex{fe.loc.x - fanEnemySize/2, fe.loc.y + fanEnemySize/2, 0.0}}
	fe.bouncers[right].bounding =
		bounding{
			vertex{fe.loc.x + fanEnemySize/2, fe.loc.y - fanEnemySize/2, 0.0},
			vertex{fe.loc.x + fanEnemySize/2 + fanEnemyMoveBoxSize, fe.loc.y + fanEnemySize/2, 0.0}}
}

func (fe *fanEnemy) collisions() []collision {
	return []collision{
		fe.harmful,
		fe.vulnerability,
		*fe.bouncers[up],
		*fe.bouncers[down],
		*fe.bouncers[left],
		*fe.bouncers[right]}
}

func (fe *fanEnemy) location() vertex {
	return fe.loc
}

func (fe *fanEnemy) mass() float64 {
	return fanEnemyMass
}

func (fe *fanEnemy) deletable() bool {
	return fe.hp <= 0
}
