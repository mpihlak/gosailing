package main

import (
	"fmt"
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
	startTime  time.Time
	paused     bool
	delayMs    int
	started    bool
	finished   bool
}

func NewSailRace() *SailRace {
	return &SailRace{
		raceCourse: gosailing.NewRaceCourse(markLocationX, markLocationY, windDirection),
		boat:       gosailing.NewBoat(maxWidth/2, 0, windDirection),
		wind:       gosailing.NewWindShifter(windDirection, 5.0, 10),
		track:      gosailing.NewTrackPlotter(),
		delayMs:    100,
	}
}

func (sr *SailRace) StartRace() {
	sr.started = true
	sr.startTime = time.Now()
}

func (sr *SailRace) Draw(win *opengl.Window) {
	win.Clear(colornames.Lightblue)
	sr.boat.Drawable().Draw(win)
	sr.raceCourse.Drawable().Draw(win)
	sr.track.Drawable().Draw(win)
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
		if win.Pressed(pixel.KeyQ) || win.Pressed(pixel.KeyEscape) {
			break
		}
		if win.Pressed(pixel.Key1) {
			sailRace.delayMs = max(0, sailRace.delayMs-10)
			fmt.Printf("Faster, delay = %v\n", sailRace.delayMs)
		}
		if win.Pressed(pixel.Key2) {
			sailRace.delayMs += 10
			fmt.Printf("Slower, delay = %v\n", sailRace.delayMs)
		}

		if win.Pressed(pixel.KeySpace) {
			if sailRace.started {
				sailRace.paused = !sailRace.paused
			} else {
				sailRace.StartRace()
			}
			time.Sleep(100 * time.Millisecond)
		}

		if sailRace.started {
			if win.Pressed(pixel.KeyT) {
				if time.Since(sailRace.lastTack) > 500*time.Millisecond {
					sailRace.boat.Tack()
					sailRace.lastTack = time.Now()
				}
			}
			if win.Pressed(pixel.KeyR) {
				sailRace = NewSailRace()
			}
		}

		sailRace.Draw(win)
		win.Update()

		currentBoatX, currentBoatY := sailRace.boat.GetXY()
		if sailRace.started && !sailRace.paused && currentBoatY <= markLocationY {
			sailRace.track.PlotLocation(sailRace.boat.GetXY())
			sailRace.boat.Advance()
			windDirection := sailRace.wind.GetWindDirection()
			sailRace.boat.SetWindDirection(windDirection)
			sailRace.raceCourse.SetWindDirection(windDirection)
		} else if currentBoatY > markLocationY {
			if !sailRace.finished {
				if currentBoatX < sailRace.raceCourse.MarkX {
					fmt.Printf("Wrong side of the mark!\n")
				}
				fmt.Printf("Sailed distance: %.2f pixels\n", sailRace.boat.GetSailedDistance())
				fmt.Printf("Sailed time: %v seconds\n", time.Since(sailRace.startTime).Seconds())
				sailRace.finished = true
			}
		}

		time.Sleep(time.Duration(sailRace.delayMs) * time.Millisecond)
	}
}

func main() {
	opengl.Run(run)
}
