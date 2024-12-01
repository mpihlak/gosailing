package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.io/mpihlak/gosailing"
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
	steps      int
	paused     bool
	delayMs    int
	started    bool
	finished   bool
	race       *imdraw.IMDraw
}

func NewSailRace() *SailRace {
	return &SailRace{
		raceCourse: gosailing.NewRaceCourse(markLocationX, markLocationY, windDirection),
		boat:       gosailing.NewBoat(maxWidth/2, 0, windDirection),
		wind:       gosailing.NewWindShifter(windDirection, 10.0, 10),
		track:      gosailing.NewTrackPlotter(),
		delayMs:    50,
	}
}

func (sr *SailRace) StartRace() {
	sr.started = true
}

func (sr *SailRace) Draw(win *opengl.Window) {
	sr.boat.Drawable().Draw(win)
	sr.raceCourse.Drawable().Draw(win)
	sr.track.Drawable().Draw(win)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	if !sr.started {
		basicTxt := text.New(pixel.V(50, 100), basicAtlas)
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, "Press SPACE to start, 'q' to quit")
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
	}
	if sr.paused {
		basicTxt := text.New(pixel.V(50, 100), basicAtlas)
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, "Paused, press SPACE to unpause")
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
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
		if win.Pressed(pixel.KeyQ) || win.Pressed(pixel.KeyEscape) {
			break
		}
		if win.Pressed(pixel.Key1) {
			sailRace.delayMs = max(0, sailRace.delayMs-10)
		}
		if win.Pressed(pixel.Key2) {
			sailRace.delayMs += 10
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

		win.Clear(colornames.Lightblue)

		currentBoatX, currentBoatY := sailRace.boat.GetXY()
		if sailRace.started && !sailRace.paused && currentBoatY <= markLocationY {
			sailRace.track.PlotLocation(sailRace.boat.GetXY())
			sailRace.boat.Advance()
			windDirection := sailRace.wind.GetWindDirection()
			sailRace.boat.SetWindDirection(windDirection)
			sailRace.raceCourse.SetWindDirection(windDirection)
			sailRace.steps++
		} else if currentBoatY > markLocationY {
			sailRace.finished = true
			basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			basicTxt := text.New(pixel.V(50, maxHeight-50), basicAtlas)
			basicTxt.Color = colornames.Black

			if currentBoatX < sailRace.raceCourse.MarkX {
				fmt.Fprintln(basicTxt, "Wrong side of the mark!")
			}

			extraDistance := math.Hypot(currentBoatX-markLocationX, currentBoatY-markLocationY)
			fmt.Fprintf(basicTxt, "Sailed distance: %.2f pixels (+extra %.2f to mark)\n", sailRace.boat.GetSailedDistance(), extraDistance)
			fmt.Fprintf(basicTxt, "Sailed time: %d units\n", sailRace.steps)

			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
		}

		sailRace.Draw(win)
		win.Update()
		time.Sleep(time.Duration(sailRace.delayMs) * time.Millisecond)
	}
}

func main() {
	opengl.Run(run)
}
