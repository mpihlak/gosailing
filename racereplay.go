package gosailing

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.io/mpihlak/gosailing/datasource"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type RaceReplay struct {
	raceCourse *RaceCourse
	boat       *Boat
	track      *TrackPlotter
	replayData []replayDataPoint
	xOffset    float64
	yOffset    float64
	currentPos int
	delayMs    int
	paused     bool
	started    bool
	finished   bool
	laylines   bool
	race       *imdraw.IMDraw
}

type replayDataPoint struct {
	datasource.NavigationDataPoint
	x float64
	y float64
}

func NewRaceReplay(markLat, markLng, maxWidth, maxHeight, zoomLevel float64, replayData *datasource.ReplayNavigationDataProvider) (*RaceReplay, error) {
	navDataPoints := replayData.GetAllPoints()
	if len(navDataPoints) == 0 {
		return nil, errors.New("no navigation data points found")
	}

	medianWind := datasource.MedianWindDirection(navDataPoints)

	markX, markY := LatLngToScreen(markLat, markLng, zoomLevel)

	// Rotate the replay points to the median wind direction
	replayDataPoints := make([]replayDataPoint, len(navDataPoints))
	var minX, minY, maxX, maxY float64
	for i, p := range navDataPoints {
		x, y := LatLngToScreen(p.Latitude, p.Longitude, zoomLevel)
		x, y = RotatePoint(x, y, markX, markY, -medianWind)

		replayDataPoints[i] = replayDataPoint{NavigationDataPoint: p, x: x, y: y}
		replayDataPoints[i].CourseOverGround -= medianWind
		replayDataPoints[i].TrueWindDirection -= medianWind

		if i == 0 || x < minX {
			minX = x
		}
		if i == 0 || y < minY {
			minY = y
		}
		if i == 0 || x > maxX {
			maxX = x
		}
		if i == 0 || y > maxY {
			maxY = y
		}
	}

	xOffset := markX - maxWidth/2
	yOffset := minY - 50

	p := replayDataPoints[0]
	return &RaceReplay{
		raceCourse: NewRaceCourse(markX-xOffset, markY-yOffset, p.TrueWindDirection),
		replayData: replayDataPoints,
		boat:       NewBoat(p.x-xOffset, p.y-yOffset, p.TrueWindDirection),
		track:      NewTrackPlotter(markX-xOffset, markY-yOffset),
		delayMs:    50,
		laylines:   true,
		xOffset:    xOffset,
		yOffset:    yOffset,
	}, nil
}

func (rr *RaceReplay) StartReplay() {
	rr.currentPos = 0
	rr.started = true
	rr.finished = false
	rr.track.Clear()
}

func (rr *RaceReplay) IsFinished() bool {
	return rr.finished
}

func (rr *RaceReplay) IncreaseSpeed() {
	rr.delayMs = max(0, rr.delayMs-10)
}

func (rr *RaceReplay) DecreaseSpeed() {
	rr.delayMs += 10
}

func (rr *RaceReplay) TogglePause() {
	if rr.started {
		rr.paused = !rr.paused
	} else {
		rr.StartReplay()
	}
}

func (rr *RaceReplay) ToggleLaylines() {
	rr.boat.ToggleLaylines()
	rr.raceCourse.ToggleLaylines()
}

func (rr *RaceReplay) ToggleWindDirection() {
	rr.raceCourse.ToggleWindDirection()
}

func (rr *RaceReplay) Throttle() {
	time.Sleep(time.Duration(rr.delayMs) * time.Millisecond)
}

func (rr *RaceReplay) Update(win *opengl.Window) {
	windowBounds := win.Bounds()
	topLeftY := windowBounds.H()

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	if !rr.started {
		textX := windowBounds.Center().X - 150
		textY := windowBounds.Center().Y + 100

		lines := []string{
			"Press SPACE to start or pause",
			"'q' quits",
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
	} else {
		navData := rr.replayData[rr.currentPos]
		if !rr.paused {
			if rr.currentPos < len(rr.replayData)-1 {
				rr.currentPos++
			} else {
				rr.finished = true
			}
		}

		rr.boat.SetLocation(navData.x-rr.xOffset, navData.y-rr.yOffset, navData.CourseOverGround, navData.TrueWindDirection)
		rr.track.PlotLocation(rr.boat.GetXY())
		rr.raceCourse.SetWindDirection(navData.TrueWindDirection)

		basicTxt := text.New(pixel.V(10, topLeftY-25), basicAtlas)
		basicTxt.Color = colornames.Black

		currentBoatX, currentBoatY := rr.boat.GetXY()
		distanceToMark := math.Hypot(currentBoatX-rr.raceCourse.MarkX, currentBoatY-rr.raceCourse.MarkY)
		fmt.Fprintf(basicTxt, "Sailed distance:  %.2f\n", rr.boat.GetSailedDistance())
		fmt.Fprintf(basicTxt, "Distance to mark: %.2f\n", distanceToMark)

		twd := -rr.raceCourse.windDirection
		if twd < 0 {
			twd = 360 + twd
		}
		fmt.Fprintf(basicTxt, "TWD: %03.0f\n", twd)
		hdg := -rr.boat.heading
		if hdg < 0 {
			hdg += 360
		}
		fmt.Fprintf(basicTxt, "HDG: %03.0f\n", hdg)
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))

		if rr.paused {
			textX := windowBounds.Center().X
			textY := windowBounds.Center().Y
			pauseBanner := "Paused, press SPACE to unpause"

			basicTxt := text.New(pixel.V(textX, textY), basicAtlas)
			basicTxt.Dot.X -= basicTxt.BoundsOf(pauseBanner).W() / 2
			basicTxt.Color = colornames.Black
			fmt.Fprintln(basicTxt, pauseBanner)
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
		}

		if rr.finished {
			textX := windowBounds.Center().X - 100
			textY := windowBounds.Center().Y

			basicTxt := text.New(pixel.V(textX, textY), basicAtlas)

			basicTxt.Color = colornames.Darkblue
			fmt.Fprintf(basicTxt, "TOTAL DISTANCE: %.2f\n", rr.boat.GetSailedDistance()+distanceToMark)
			if currentBoatX < rr.raceCourse.MarkX {
				basicTxt.Color = colornames.Red
				fmt.Fprintln(basicTxt, "Wrong side of the mark!")
			}
			basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
		}

		rr.boat.Drawable().Draw(win)
		rr.raceCourse.Drawable().Draw(win)
		rr.track.Drawable().Draw(win)
	}
}
