package main

import (
	"time"
)

var coolGuyCol int
var otherCoolGuy int

const (
	spawnClearerRadius = 300
	spawnDelay         = 3 * time.Second
)

type playerRespawner struct {
	uidGenerator

	spawnClearers []collision
	players       []*spaceShipListener
	defeatTimes   map[*spaceShipListener]time.Time
}

func newPlayerRespawner(players []*spaceShipListener) *playerRespawner {
	pr := new(playerRespawner)

	otherCoolGuy = pr.uid()

	pr.spawnClearers = []collision{}
	pr.players = players
	pr.defeatTimes = make(map[*spaceShipListener]time.Time)

	return pr
}

func (pr *playerRespawner) tick() []entity {
	respawns := []entity{}

	// Spawn clearers should only exist for one frame
	pr.spawnClearers = []collision{}

	// Record the time a player is defeated
	for _, player := range pr.players {
		if _, recorded := pr.defeatTimes[player]; player.deletable() && !recorded {
			pr.defeatTimes[player] = time.Now()
		}
	}

	// Respawn the player after a delay
	for player, defeatTime := range pr.defeatTimes {
		if time.Since(defeatTime) >= spawnDelay {
			player.respawn()
			pr.spawnClearers = append(pr.spawnClearers, pr.newSpawnClearer())
			respawns = append(respawns, player)
			delete(pr.defeatTimes, player)
		}
	}

	return respawns
}

// newSpawnClearer creates a hazardous collision meant to wipe out all enemies
// around the player's spawn point.
func (pr *playerRespawner) newSpawnClearer() collision {
	return collision{
		bounding: bounding{
			vertex{-spawnClearerRadius, -spawnClearerRadius, 0.0},
			vertex{spawnClearerRadius, spawnClearerRadius, 0.0}},
		alliance: friendly,
		typ:      deadly}
}

func (pr *playerRespawner) collisions() []collision {
	cols := []collision{}

	cols = append(cols, pr.spawnClearers...)

	return cols
}

func (pr *playerRespawner) collision(yours, other collision) {
}

func (pr *playerRespawner) draw() {
}

func (pr *playerRespawner) location() vertex {
	return vertex{0, 0, 0}
}

func (pr *playerRespawner) deletable() bool {
	return false
}

func (pr *playerRespawner) mass() float64 {
	return 0.0
}

func init() {
	var pe physicalEntity = &playerRespawner{}
	_ = pe
}
