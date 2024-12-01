package main

import (
	"flag"
	"fmt"
	"time"

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

var (
	windData = flag.String("windData", "", "Wind data file")
)

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Go Sailing!",
		Bounds: pixel.R(0, 0, maxWidth, maxHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	newSailRace := func() *gosailing.SailRace {
		var windShifter gosailing.WindShifter
		if *windData != "" {
			fmt.Printf("Using wind data from %v\n", *windData)
			windShifter = gosailing.NewReplayShifter(*windData)
		} else {
			fmt.Printf("Generating oscillating wind")
			windShifter = gosailing.NewOscillatingWindShifter(windDirection, 10.0, 10)
		}

		return gosailing.NewSailRace(
			markLocationX, markLocationY,
			boatLocationX, boatLocationY,
			windShifter,
		)
	}

	sailRace := newSailRace()

	// Throttle the keyboard to avoid registering unintended repeated keypresses
	lastKeyPressed := make(map[pixel.Button]time.Time)
	keyPressed := func(k pixel.Button) bool {
		if win.Pressed(k) && time.Since(lastKeyPressed[k]) > 200*time.Millisecond {
			lastKeyPressed[k] = time.Now()
			return true
		}
		return false
	}

	for !win.Closed() {
		if keyPressed(pixel.KeyQ) || win.Pressed(pixel.KeyEscape) {
			break
		}
		if keyPressed(pixel.Key1) {
			sailRace.IncreaseSpeed()
		}
		if keyPressed(pixel.Key2) {
			sailRace.DecreaseSpeed()
		}
		if keyPressed(pixel.KeySpace) {
			sailRace.TogglePause()
		}
		if keyPressed(pixel.KeyT) {
			sailRace.TackBoat()
		}
		if keyPressed(pixel.KeyR) {
			sailRace = newSailRace()
		}
		if keyPressed(pixel.KeyL) {
			sailRace.ToggleLaylines()
		}
		if keyPressed(pixel.KeyW) {
			sailRace.ToggleWindDirection()
		}

		win.Clear(colornames.Lightblue)
		sailRace.Update(win)
		win.Update()

		sailRace.Throttle()
	}
}

func main() {
	flag.Parse()
	opengl.Run(run)
}
