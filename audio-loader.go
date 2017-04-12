package main

func loadAudio() error {
	err := loadMusic("resources/music/title-loop.ogg", titleMusic)
	if err != nil {
		return err
	}

	err = loadMusic("resources/music/gameplay-loop.ogg", gameplayMusic)
	if err != nil {
		return err
	}

	return nil
}
