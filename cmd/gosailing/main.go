package main

import (
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.io/mpihlak/gosailing"
	"golang.org/x/image/colornames"
)

const (
	maxWidth  = 1024.0
	maxHeight = 768.0

	windDirection = 10.0
	markLocationX = maxWidth / 2
	markLocationY = maxHeight - 50
)

type SailRace struct {
	raceCourse *gosailing.RaceCourse
	boat       *gosailing.Boat
	wind       *gosailing.WindShifter
	track      *gosailing.TrackPlotter
	lastTack   time.Time
	paused     bool
}

func NewSailRace() *SailRace {
	return &SailRace{
		raceCourse: gosailing.NewRaceCourse(markLocationX, markLocationY, windDirection),
		boat:       gosailing.NewBoat(maxWidth/2, 0, windDirection),
		wind:       gosailing.NewWindShifter(windDirection, 15.0, 10),
		track:      gosailing.NewTrackPlotter(),
	}
}

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

	sailRace := NewSailRace()
	for !win.Closed() {
		if win.Pressed(pixel.KeySpace) {
			sailRace.paused = !sailRace.paused
		}

		if win.Pressed(pixel.KeyQ) || win.Pressed(pixel.KeyEscape) {
			break
		}

		if win.Pressed(pixel.KeyT) {
			if time.Since(sailRace.lastTack) > 500*time.Millisecond {
				sailRace.boat.Tack()
				sailRace.lastTack = time.Now()
			}
		}

		if win.Pressed(pixel.KeyR) {
			sailRace = NewSailRace()
		}

		win.Clear(colornames.Lightblue)
		sailRace.boat.Drawable().Draw(win)
		sailRace.raceCourse.Drawable().Draw(win)
		sailRace.track.Drawable().Draw(win)
		win.Update()

		_, currentBoatY := sailRace.boat.GetXY()
		if !sailRace.paused && currentBoatY <= markLocationY {
			sailRace.track.PlotLocation(sailRace.boat.GetXY())
			sailRace.boat.Advance()
			windDirection := sailRace.wind.GetWindDirection()
			sailRace.boat.SetWindDirection(windDirection)
			sailRace.raceCourse.SetWindDirection(windDirection)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	opengl.Run(run)
}
