package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	spawnDistFromEntities = 150
	spawnInterval         = time.Second
	spawnBorder           = 40.0
)

type enemySpawner struct {
	uidGenerator
	spawnTicker <-chan time.Time
}

func newEnemySpawner() *enemySpawner {
	es := new(enemySpawner)

	es.spawnTicker = time.Tick(spawnInterval)

	return es
}

func (es *enemySpawner) tick() []entity {
	var newEnts []entity

	select {
	case <-es.spawnTicker:
		var done bool
		var loc vertex
		var tries int
		for !done {
			loc = vertex{
				(rand.Float64()*2.0 - 1.0) * ((playAreaWidth - spawnBorder) / 2.0),
				(rand.Float64()*2.0 - 1.0) * ((playAreaHeight - spawnBorder) / 2.0),
				0.0}
			done = true

			for _, pe := range currentReg.physicals {
				if pe.location().distance(loc) < spawnDistFromEntities {
					done = false
				}
			}
			tries++
		}

		newEnts = append(newEnts, newFanEnemy(loc))
		fmt.Println("Spawned an enemy, took", tries, "tries")
	default:
	}

	return newEnts
}

func (es *enemySpawner) deletable() bool {
	return false
}
