package main

var mainCamera camera

type camera struct {
	uidGenerator

	loc vertex
}

func (c *camera) tick() []entity {
	var sumLoc vertex

	inFramers := currentReg.inFramers

	if len(inFramers) > 0 {
		min := inFramers[0].location()
		max := inFramers[0].location()

		// Calculate the average, min, and max locations of entities
		for _, ent := range inFramers {
			loc := ent.location()
			sumLoc.x += loc.x
			sumLoc.y += loc.y

			if loc.x < min.x {
				min.x = loc.x
			}
			if loc.x > max.x {
				max.x = loc.x
			}
			if loc.y < min.y {
				min.y = loc.y
			}
			if loc.y > max.y {
				max.y = loc.y
			}
		}

		// Zoom out the camera so all entities can be seen at once
		var zScaling vertex
		if max.x-min.x > float64(mainWindow.width)/2.0 {
			zScaling.x = (max.x - min.x) / (float64(mainWindow.width) / 2.0)
		} else {
			zScaling.x = 1
		}
		if max.y-min.y < float64(mainWindow.height)/2.0 {
			zScaling.y = (max.y - min.y) / (float64(mainWindow.height) / 2.0)
		} else {
			zScaling.y = 1
		}

		largestZScale := zScaling.x
		if zScaling.y > largestZScale {
			largestZScale = zScaling.y
		}

		c.loc = vertex{
			x: sumLoc.x / float64(len(inFramers)),
			y: sumLoc.y / float64(len(inFramers)),
			z: largestZScale}
	}

	return []entity{}
}

func (c *camera) deletable() bool {
	return false
}
