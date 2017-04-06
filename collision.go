package main

type collisionType int

const (
	bouncer collisionType = iota
	harmful
	vulnerable
)

type alliance int

const (
	friendly alliance = iota
	unfriendly
	unaligned
)

type collision struct {
	uidGenerator

	bounding bounding
	alliance alliance
	typ      collisionType
}

func collides(c1, c2 collision) bool {
	if c1.bounding.b.x < c2.bounding.a.x {
		return false
	}
	if c1.bounding.a.x > c2.bounding.b.x {
		return false
	}
	if c1.bounding.b.y < c2.bounding.a.y {
		return false
	}
	if c1.bounding.a.y > c2.bounding.b.y {
		return false
	}

	return true
}
