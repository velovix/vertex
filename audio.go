package main

import (
	"time"

	mix "github.com/veandco/go-sdl2/mix"
)

type music int

const (
	titleMusic music = iota
	gameplayMusic
)

var (
	loadedMusic map[music]*mix.Music
	musicMuted  bool
)

func init() {
	err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}

	loadedMusic = make(map[music]*mix.Music)
}

func loadMusic(filename string, mus music) error {
	sdlMus, err := mix.LoadMUS(filename)
	if err != nil {
		return err
	}

	loadedMusic[mus] = sdlMus

	return nil
}

func playMusic(m music, fade time.Duration) {
	if !musicMuted {
		loadedMusic[m].FadeIn(-1, int(fade/time.Millisecond))
	}
}

func stopMusic() {
	mix.HaltMusic()
}

func musicPlaying() bool {
	return mix.PlayingMusic()
}
