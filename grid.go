package main

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

const (
	// gridMargin controls the distance between grid points.
	gridMargin = 30
	// gridDestructionForce controls the amount of ripple created from entity
	// destruction.
	gridDestructionForce = 50
	// gridMoveForce controls the amount of ripple created from entity
	// movement.
	gridMoveForce = 10
)

// grid provides a visual grid with water rippling effects.
type grid struct {
	uidGenerator

	points [][]gridPoint
}

func newGrid(width, height float64) *grid {
	g := &grid{}

	// Make a new 2D array of points
	g.points = make([][]gridPoint, int(math.Ceil(width/gridMargin)))
	for i := range g.points {
		g.points[i] = make([]gridPoint, int(math.Ceil(height/gridMargin)))
	}

	// Place all points
	for x := range g.points {
		for y := range g.points[x] {
			g.points[x][y] = newGridPoint(vertex{
				(-width / 2.0) + float64(x)*gridMargin,
				(-height / 2.0) + float64(y)*gridMargin,
				0.0})
		}
	}

	return g
}

func (g *grid) tick() []entity {
	// Update every point
	for x := range g.points {
		for y := range g.points[x] {
			g.points[x][y].tick()
		}
	}

	return []entity{}
}

func (g *grid) draw() {
	gl.PushMatrix()

	// Give the grid some generic off white material
	matAmbient := []float32{0.6, 0.6, 0.6, 1.0}
	matDiffuse := []float32{0.6, 0.6, 0.6, 1.0}
	matSpecular := []float32{0.0, 0.0, 0.0, 1.0}

	gl.Materialfv(gl.FRONT, gl.AMBIENT, fPtr(matAmbient))
	gl.Materialfv(gl.FRONT, gl.DIFFUSE, fPtr(matDiffuse))
	gl.Materialfv(gl.FRONT, gl.SPECULAR, fPtr(matSpecular))
	gl.Materialf(gl.FRONT, gl.SHININESS, 0)

	// Draw lines between every point
	gl.Begin(gl.LINES)
	for x := 0; x < len(g.points); x++ {
		for y := 0; y < len(g.points[x]); y++ {
			if x+1 < len(g.points) {
				gl.Normal3d(0, 0, 1)
				gl.Vertex3d(
					g.points[x][y].loc.x,
					g.points[x][y].loc.y,
					gridPointZVal)
				gl.Normal3d(0, 0, 1)
				gl.Vertex3d(
					g.points[x+1][y].loc.x,
					g.points[x+1][y].loc.y,
					gridPointZVal)
			}
			if y+1 < len(g.points[x]) {
				gl.Normal3d(0, 0, 1)
				gl.Vertex3d(
					g.points[x][y].loc.x,
					g.points[x][y].loc.y,
					gridPointZVal)
				gl.Normal3d(0, 0, 1)
				gl.Vertex3d(
					g.points[x][y+1].loc.x,
					g.points[x][y+1].loc.y,
					gridPointZVal)
			}
		}
	}
	gl.End()

	gl.PopMatrix()
}

func (g *grid) entityMove(e physicalEntity, loc vertex) {
	// Create a ripple around moving enemies
	// TODO(velovix): Refactor this into a method
	for x := range g.points {
		for y := range g.points[x] {
			dist := g.points[x][y].loc.distance(loc)
			// TODO(velovix): assign a constant to this 25 value
			magnitude := (gridMoveForce / (dist / 25.0)) * e.mass()
			if magnitude > gridMoveForce {
				magnitude = gridMoveForce
			}

			direction := normalize(loc.subtract(g.points[x][y].loc))

			g.points[x][y].applyForce(vertex{
				-direction.x * magnitude,
				-direction.y * magnitude,
				0})
		}
	}
}

func (g *grid) entityDestruction(e entity) {
	// Create a ripple around destroyed enemies
	// TODO(velovix): Refactor this into a method
	if pe, ok := e.(physicalEntity); ok {
		for x := range g.points {
			for y := range g.points[x] {
				dist := g.points[x][y].loc.distance(pe.location())
				baseForce := gridDestructionForce * pe.mass()
				// TODO(velovix): assign a constnant to this 50 value
				magnitude := baseForce / (dist / (50.0 * pe.mass()))
				if magnitude > gridDestructionForce*pe.mass() {
					magnitude = baseForce
				}

				direction := normalize(pe.location().subtract(g.points[x][y].loc))

				g.points[x][y].applyForce(vertex{
					-direction.x * magnitude,
					-direction.y * magnitude,
					0})
			}
		}
	}
}

func (g *grid) location() vertex {
	return vertex{0, 0, 0}
}

func (g *grid) collisions() []collision {
	return []collision{}
}

func (g *grid) collision(yours, other collision) {
}

func (g *grid) mass() float64 {
	return 0
}

func (g *grid) deletable() bool {
	return false
}

const (
	// gridPointMass is an arbitrary mass value.
	gridPointMass = 3
	// gridPointSpringConst controls how quickly a point snaps back to the
	// direction of its original location.
	gridPointSpringConst = 1.0
	// gridPointSprintDampening dampens the point's movement so it eventually
	// stops.
	gridPointSpringDampening = 0.9

	// gridPointZVal is the Z value of all grid points.
	gridPointZVal = -3
)

// gridPoint is a single point on the grid.
type gridPoint struct {
	loc, origLoc, accel vertex
}

func newGridPoint(loc vertex) gridPoint {
	return gridPoint{
		loc:     loc,
		origLoc: loc}
}

func (gp *gridPoint) tick() {
	// Apply the current acceleration
	gp.loc.x += gp.accel.x * mainWindow.delta
	gp.loc.y += gp.accel.y * mainWindow.delta

	// Affect acceleration using Hooke's law
	gp.accel.x += ((gridPointSpringConst * (gp.origLoc.x - gp.loc.x)) / gridPointMass)
	gp.accel.y += ((gridPointSpringConst * (gp.origLoc.y - gp.loc.y)) / gridPointMass)

	// Dampen the acceleration
	gp.accel.x *= (gridPointSpringDampening)
	gp.accel.y *= (gridPointSpringDampening)
}

func (gp *gridPoint) applyForce(force vertex) {
	gp.accel.x += force.x / gridPointMass
	gp.accel.y += force.y / gridPointMass
}
