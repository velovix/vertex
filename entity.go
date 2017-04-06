package main

import "github.com/go-gl/glfw/v3.0/glfw"

// entity represents a non-descript object that does operations on every frame.
type entity interface {
	// uid returns a unique identifier.
	uid() int

	// tick is called once per frame. Returned entities are added to the entity
	// registry.
	tick() []entity

	// deletable returns true if the entity should be deleted.
	deletable() bool
}

// physicalEntity represents an entity that takes on some physical, usually
// viewiable presence.
type physicalEntity interface {
	entity

	// draw is called once per frame. Implementations should draw to the
	// screen.
	draw()

	// location returns the entity's location.
	location() vertex

	// collisions returns all owned collisions of the entity.
	collisions() []collision

	// collision is called when a collision happens against this entity.
	collision(yours, other collision)

	// mass returns some arbitrary value giving callers an idea of how "heavy"
	// the entity is.
	mass() float64
}

// inFrameEntity is an entity that should always be in frame. A player entity,
// for instance. Entities need to physically exist to effect camera location.
type inFrameEntity interface {
	physicalEntity

	// inFrame means nothing, but should be implemented by those who want to be
	// always in frame.
	inFrame()
}

// inputEntity represents an entity that needs to know about user input.
type inputEntity interface {
	entity

	// joystickAxis is called when a joystick axis value changes.
	joystickAxis(joystick glfw.Joystick, axes []float32)

	// joystickButton is called when a joystick button value changes.
	joystickButton(joystick glfw.Joystick, buttons []byte)
}

// snoopingEntity represents an entity that needs to respond to the changing
// state of other entities.
type snoopingEntity interface {
	entity

	// entityDestruction is called right before an entity is destroyed.
	entityDestruction(e entity)

	// entityMove is called when an entity moves around.
	entityMove(e physicalEntity, loc vertex)
}
