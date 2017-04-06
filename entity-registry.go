package main

import (
	"github.com/go-gl/glfw/v3.0/glfw"
)

// currentReg holds the in-focus entity registry.
var currentReg *entityRegistry

// entityRegistry is an owner of entities. It does operations on these entities
// when events happen, depending on what entity interfaces the entity
// implements.
type entityRegistry struct {
	entities  []entity
	physicals []physicalEntity
	inputters []inputEntity
	snoopers  []snoopingEntity
	inFramers []inFrameEntity
}

// addEntity adds an entity to the registry.
func (er *entityRegistry) addEntity(e entity) {
	er.entities = append(er.entities, e)

	if physical, ok := e.(physicalEntity); ok {
		er.physicals = append(er.physicals, physical)
	}
	if inputter, ok := e.(inputEntity); ok {
		er.inputters = append(er.inputters, inputter)
	}
	if snooper, ok := e.(snoopingEntity); ok {
		er.snoopers = append(er.snoopers, snooper)
	}
	if inFramer, ok := e.(inFrameEntity); ok {
		er.inFramers = append(er.inFramers, inFramer)
	}
}

// draw draws every physical entity.
func (er *entityRegistry) draw() {
	for _, e := range er.physicals {
		e.draw()
	}
}

// tick runs the tick operations of every entity.
func (er *entityRegistry) tick() {
	for _, e := range er.entities {
		var newEnts []entity

		if pe, ok := e.(physicalEntity); ok {
			// Find out if the entity moved this tick
			oldLoc := pe.location()
			newEnts = pe.tick()
			newLoc := pe.location()
			if newLoc != oldLoc {
				er.moveEvent(pe, newLoc)
			}
		} else {
			newEnts = e.tick()
		}

		// Add any new entities to registry
		for _, newEnt := range newEnts {
			er.addEntity(newEnt)
		}

		// Delete the entity if necessary
		if e.deletable() {
			er.deleteEvent(e)
			er.del(e)
		}
	}

	// Test for collisions between entities
	for i := 0; i < len(er.physicals); i++ {
		for j := i + 1; j < len(er.physicals); j++ {
			er.interact(er.physicals[i], er.physicals[j])
		}
	}
}

// interact has the two given entities' collisions interact with each other.
func (er *entityRegistry) interact(pe1 physicalEntity, pe2 physicalEntity) {
	cols1 := pe1.collisions()
	cols2 := pe2.collisions()

	for _, c1 := range cols1 {
		for _, c2 := range cols2 {
			if collides(c1, c2) {
				pe1.collision(c1, c2)
				pe2.collision(c2, c1)
			}
		}
	}
}

// moveEvent notifies all snooping entities that the given entity has moved.
func (er *entityRegistry) moveEvent(e physicalEntity, loc vertex) {
	for _, se := range er.snoopers {
		se.entityMove(e, loc)
	}
}

// deleteEvent notifies all snooping entities that the given entity has been
// deleted.
func (er *entityRegistry) deleteEvent(e entity) {
	for _, se := range er.snoopers {
		se.entityDestruction(e)
	}
}

// joystickAxisEvent notifies all input entities that a joystick axis has
// changed.
func (er *entityRegistry) joystickAxisEvent(joystick glfw.Joystick, axes []float32) {
	for _, ie := range er.inputters {
		ie.joystickAxis(joystick, axes)
	}
}

// joystickButtonEvent notifies all input entities that a joystick button has
// changed.
func (er *entityRegistry) joystickButtonEvent(joystick glfw.Joystick, buttons []byte) {
	for _, ie := range er.inputters {
		ie.joystickButton(joystick, buttons)
	}
}

// del deletes the given entity from the entity registry.
func (er *entityRegistry) del(e entity) {
	for i := 0; i < len(er.entities); i++ {
		if er.entities[i].uid() == e.uid() {
			er.entities = append(er.entities[:i], er.entities[i+1:]...)
		}
	}

	for i := 0; i < len(er.physicals); i++ {
		if er.physicals[i].uid() == e.uid() {
			er.physicals = append(er.physicals[:i], er.physicals[i+1:]...)
		}
	}

	for i := 0; i < len(er.inputters); i++ {
		if er.inputters[i].uid() == e.uid() {
			er.inputters = append(er.inputters[:i], er.inputters[i+1:]...)
		}
	}

	for i := 0; i < len(er.snoopers); i++ {
		if er.snoopers[i].uid() == e.uid() {
			er.snoopers = append(er.snoopers[:i], er.snoopers[i+1:]...)
		}
	}
}
