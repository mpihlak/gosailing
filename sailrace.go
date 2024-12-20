package gosailing

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
)

type SailRace struct {
	raceCourse *RaceCourse
	boat       *Boat
	wind       WindShifter
	track      *TrackPlotter
	lastTack   time.Time
	delayMs    int
	paused     bool
	started    bool
	finished   bool
	laylines   bool
	race       *imdraw.IMDraw
}

func NewSailRace(markLocationX, markLocationY, boatLocationX, boatLocationY float64, windShifter WindShifter) *SailRace {
	wd := windShifter.GetWindDirection()
	return &SailRace{
		raceCourse: NewRaceCourse(markLocationX, markLocationY, wd),
		boat:       NewBoat(boatLocationX, boatLocationY, wd),
		wind:       windShifter,
		track:      NewTrackPlotter(boatLocationX, boatLocationY),
		delayMs:    50,
		laylines:   true,
	}
}

func (sr *SailRace) StartRace() {
	sr.started = true
}

func (sr *SailRace) IsFinished() bool {
	return sr.finished
}

func (sr *SailRace) IncreaseSpeed() {
	sr.delayMs = max(0, sr.delayMs-10)
}

func (sr *SailRace) DecreaseSpeed() {
	sr.delayMs += 10
}

func (sr *SailRace) TogglePause() {
	if sr.started {
		sr.paused = !sr.paused
	} else {
		sr.StartRace()
	}
}

func (sr *SailRace) TackBoat() {
	if sr.started {
		if time.Since(sr.lastTack) > 500*time.Millisecond {
			sr.boat.Tack()
			sr.lastTack = time.Now()
		}
	}
}

func (sr *SailRace) ToggleLaylines() {
	sr.boat.ToggleLaylines()
	sr.raceCourse.ToggleLaylines()
}

func (sr *SailRace) ToggleWindDirection() {
	sr.raceCourse.ToggleWindDirection()
}

func (sr *SailRace) Throttle() {
	time.Sleep(time.Duration(sr.delayMs) * time.Millisecond)
}

func (sr *SailRace) Update(win *opengl.Window) {
	windowBounds := win.Bounds()
	topLeftY := windowBounds.H()

	currentBoatX, currentBoatY := sr.boat.GetXY()

	// TODO: refactor the finishing condition so that it also works on non-vertical courses

	if sr.started && !sr.paused && currentBoatY <= sr.raceCourse.MarkY {
		sr.track.PlotLocation(sr.boat.GetXY())
		sr.boat.Advance()
		windDirection := sr.wind.GetWindDirection()
		sr.boat.SetWindDirection(windDirection)
		sr.raceCourse.SetWindDirection(windDirection)
	} else if currentBoatY > sr.raceCourse.MarkY {
		sr.finished = true
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	if sr.started {
		basicTxt := text.New(pixel.V(10, topLeftY-25), basicAtlas)
		basicTxt.Color = colornames.Black

		distanceToMark := math.Hypot(currentBoatX-sr.raceCourse.MarkX, currentBoatY-sr.raceCourse.MarkY)
		fmt.Fprintf(basicTxt, "Sailed distance:  %.2f\n", sr.boat.GetSailedDistance())
		fmt.Fprintf(basicTxt, "Distance to mark: %.2f\n", distanceToMark)

		twd := -sr.raceCourse.windDirection
		if twd < 0 {
			twd = 360 + twd
		}
		fmt.Fprintf(basicTxt, "TWD: %03.0f\n", twd)
		hdg := -sr.boat.heading
		if hdg < 0 {
			hdg += 360
		}
		fmt.Fprintf(basicTxt, "HDG: %03.0f\n", hdg)
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))

		if sr.finished {
			textX := windowBounds.Center().X - 100
			textY := windowBounds.Center().Y

			basicTxt := text.New(pixel.V(textX, textY), basicAtlas)

			basicTxt.Color = colornames.Darkblue
			fmt.Fprintf(basicTxt, "TOTAL DISTANCE: %.2f\n", sr.boat.GetSailedDistance()+distanceToMark)
			if currentBoatX < sr.raceCourse.MarkX {
				basicTxt.Color = colornames.Red
				fmt.Fprintln(basicTxt, "Wrong side of the mark!")
			}
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
		}
	}

	sr.boat.Drawable().Draw(win)
	sr.raceCourse.Drawable().Draw(win)
	sr.track.Drawable().Draw(win)

	if !sr.started {
		textX := windowBounds.Center().X - 150
		textY := windowBounds.Center().Y + 100

		lines := []string{
			"Press SPACE to start or pause",
			"'q' quits",
			"'t' tacks'",
			"'r' restarts'",
			"'l' toggle laylines",
			"'w' toggle wind",
			"'1' increases speed'",
			"'2' decreases speed'",
		}

		basicTxt := text.New(pixel.V(textX, textY), basicAtlas)
		basicTxt.Color = colornames.Black
		for _, line := range lines {
			fmt.Fprintln(basicTxt, line)
		}
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
	}

	if sr.paused {
		textX := windowBounds.Center().X
		textY := windowBounds.Center().Y
		pauseBanner := "Paused, press SPACE to unpause"

		basicTxt := text.New(pixel.V(textX, textY), basicAtlas)
		basicTxt.Dot.X -= basicTxt.BoundsOf(pauseBanner).W() / 2
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, pauseBanner)
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
	}
}
