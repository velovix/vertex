package main

const (
	boundariesSize = 30
)

type boundaries struct {
	uidGenerator

	cols map[direction]collision
}

func newBoundaries(width, height float64) *boundaries {
	b := new(boundaries)

	b.cols = make(map[direction]collision)

	b.cols[up] = collision{
		bounding: bounding{
			vertex{-width / 2.0, height / 2.0, 0.0},
			vertex{width / 2.0, height/2.0 + boundariesSize, 0.0}},
		alliance: unaligned,
		typ:      bouncer}
	b.cols[down] = collision{
		bounding: bounding{
			vertex{-width / 2.0, -height/2.0 - boundariesSize, 0.0},
			vertex{width / 2.0, -height/2.0 + boundariesSize, 0.0}},
		alliance: unaligned,
		typ:      bouncer}
	b.cols[left] = collision{
		bounding: bounding{
			vertex{-width/2.0 - boundariesSize, -height / 2.0, 0.0},
			vertex{-width / 2.0, height / 2.0, 0.0}},
		alliance: unaligned,
		typ:      bouncer}
	b.cols[right] = collision{
		bounding: bounding{
			vertex{width / 2.0, -height / 2.0, 0.0},
			vertex{width/2.0 + boundariesSize, height / 2.0, 0.0}},
		alliance: unaligned,
		typ:      bouncer}

	return b
}

func (b *boundaries) tick() []entity {
	return []entity{}
}

func (b *boundaries) draw() {
}

func (b *boundaries) location() vertex {
	return vertex{0, 0, 0}
}

func (b *boundaries) collisions() []collision {
	var cols []collision
	for _, col := range b.cols {
		cols = append(cols, col)
	}
	return cols
}

func (b *boundaries) collision(yours, other collision) {
}

func (b *boundaries) mass() float64 {
	return 0
}

func (b *boundaries) deletable() bool {
	return false
}
