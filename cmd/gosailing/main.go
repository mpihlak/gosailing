package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"

	"github.io/mpihlak/gosailing"
)

const (
	maxWidth  = 1024.0
	maxHeight = 768.0

	windDirection = 0.0 // 0 is north, 90 is east and so on
	markLocationX = maxWidth / 2
	markLocationY = maxHeight - 50
	boatLocationX = maxWidth / 2
	boatLocationY = 25
)

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Go sailing!",
		Bounds: pixel.R(0, 0, maxWidth, maxHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	newSailRace := func() *gosailing.SailRace {
		return gosailing.NewSailRace(
			markLocationX, markLocationY,
			boatLocationX, boatLocationY,
			windDirection,
		)
	}

	sailRace := newSailRace()

	for !win.Closed() {
		if win.Pressed(pixel.KeyQ) || win.Pressed(pixel.KeyEscape) {
			break
		}
		if win.Pressed(pixel.Key1) {
			sailRace.IncreaseSpeed()
		}
		if win.Pressed(pixel.Key2) {
			sailRace.DecreaseSpeed()
		}
		if win.Pressed(pixel.KeySpace) {
			sailRace.TogglePause()
		}
		if win.Pressed(pixel.KeyT) {
			sailRace.TackBoat()
		}
		if win.Pressed(pixel.KeyR) {
			sailRace = newSailRace()
		}

		win.Clear(colornames.Lightblue)
		sailRace.Update(win)
		win.Update()

		sailRace.Throttle()
	}
}

func main() {
	opengl.Run(run)
}
