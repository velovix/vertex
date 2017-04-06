package main

func loadModels() error {
	var err error

	spaceShipModel, err = loadObjFile("resources/space-ship.obj")
	if err != nil {
		return err
	}

	fanEnemyModel, err = loadObjFile("resources/fan-enemy.obj")
	if err != nil {
		return err
	}

	bulletModel, err = loadObjFile("resources/bullet.obj")
	if err != nil {
		return err
	}

	return nil
}
