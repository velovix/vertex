package main

import "time"

type enemySpawner struct {
	uidGenerator
	lastSpawn time.Time
}

func (es *enemySpawner) tick() []entity {
	var newEnts []entity

	if time.Since(es.lastSpawn) > 6*time.Second {
		newEnts = append(newEnts, newFanEnemy(vertex{500, 500, 0}))
		newEnts = append(newEnts, newFanEnemy(vertex{-500, 500, 0}))
		newEnts = append(newEnts, newFanEnemy(vertex{500, -500, 0}))
		newEnts = append(newEnts, newFanEnemy(vertex{-500, -500, 0}))

		es.lastSpawn = time.Now()
	}

	return newEnts
}

func (es *enemySpawner) deletable() bool {
	return false
}
