package main

import "github.com/go-gl/gl/v2.1/gl"

var (
	bulletMaxHP = 1

	bulletScale = vertex{x: 10, y: 10, z: 10}
	bulletModel polyModel
)

// bullet is a small projectile shot by players or enemies.
type bullet struct {
	x float64
	y float64

	moveDir vertex
	speed   float64

	rot float64

	hp int
}

func newBullet() *bullet {
	bul := new(bullet)

	bul.hp = bulletMaxHP

	return bul
}

func (b *bullet) tick() {

}

func (b *bullet) draw() {
	gl.PushMatrix()

	gl.Scaled(bulletScale.x, bulletScale.y, bulletScale.z)
	bulletModel.draw()

	gl.PopMatrix()
}
