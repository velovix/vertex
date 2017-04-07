package main

import (
	"github.com/go-gl/gl/v2.1/gl"
)

var (
	bulletMaxHP   = 1
	bulletMass    = 0.5
	bulletColSize = 20.0

	bulletScale = vertex{x: 10, y: 10, z: 10}
	bulletModel polyModel
)

// bullet is a small projectile shot by players or enemies.
type bullet struct {
	uidGenerator

	loc vertex

	moveDir vertex
	speed   float64

	rot float64

	alliance alliance
	harmful  collision

	hp int
}

// bulletPosData is positional data for a bullet, used to initialize.
type bulletPosData struct {
	loc     vertex
	rot     float64
	moveDir vertex
	speed   float64
}

func newBullet(data bulletPosData, alliance alliance) *bullet {
	bul := new(bullet)

	bul.hp = bulletMaxHP

	bul.loc = data.loc
	bul.rot = data.rot
	bul.moveDir = data.moveDir
	bul.speed = data.speed
	bul.alliance = alliance

	return bul
}

func (b *bullet) tick() []entity {
	b.loc.x += (b.speed * mainWindow.delta) * b.moveDir.x
	b.loc.y += (b.speed * mainWindow.delta) * b.moveDir.y

	b.updateCols()

	return []entity{}
}

func (b *bullet) draw() {
	gl.PushMatrix()

	rotateOnPoint(b.rot, b.loc)
	gl.Translated(b.loc.x, b.loc.y, b.loc.z)
	gl.Scaled(bulletScale.x, bulletScale.y, bulletScale.z)
	bulletModel.draw()

	gl.PopMatrix()
}

func (b *bullet) collision(yours, other collision) {
	switch other.typ {
	case bouncer:
		if other.alliance == unaligned {
			// We bumped into a wall or something
			b.hp--
		}
	case vulnerable:
		if other.alliance != b.alliance {
			// We probably hit an enemy
			b.hp--
		}
	}
}

func (b *bullet) updateCols() {
	b.harmful.alliance = b.alliance
	b.harmful.typ = harmful
	b.harmful.bounding = bounding{
		vertex{b.loc.x - bulletColSize/2.0, b.loc.y - bulletColSize/2.0, 0.0},
		vertex{b.loc.x + bulletColSize/2.0, b.loc.y + bulletColSize/2.0, 0.0}}
}

func (b *bullet) collisions() []collision {
	return []collision{
		b.harmful}
}

func (b *bullet) location() vertex {
	return b.loc
}

func (b *bullet) deletable() bool {
	return b.hp <= 0
}

func (b *bullet) mass() float64 {
	return bulletMass
}
