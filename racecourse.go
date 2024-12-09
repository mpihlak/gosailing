package gosailing

import (
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type RaceCourse struct {
	MarkX             float64
	MarkY             float64
	windDirection     float64
	laylines          bool
	showWindDirection bool
	course            *imdraw.IMDraw
}

// NewRaceCourse creates a new race course with the given mark location and wind direction
func NewRaceCourse(x, y, windDirection float64) *RaceCourse {
	course := imdraw.New(nil)

	return &RaceCourse{
		MarkX:             x,
		MarkY:             y,
		windDirection:     windDirection,
		course:            course,
		laylines:          true,
		showWindDirection: true,
	}
}

// SetWindDirection sets the current wind direction in degrees
// North is 0 and is straight up.
func (rc *RaceCourse) SetWindDirection(direction float64) {
	rc.windDirection = direction
}

func (rc *RaceCourse) ToggleLaylines() {
	rc.laylines = !rc.laylines
}

func (rc *RaceCourse) ToggleWindDirection() {
	rc.showWindDirection = !rc.showWindDirection
}

func (rc *RaceCourse) Drawable() *imdraw.IMDraw {
	rc.course.Clear()

	DrawFlag(rc.course, rc.MarkX, rc.MarkY)

	if rc.laylines {
		LayLine(rc.course, rc.MarkX, rc.MarkY, -TackAngle+rc.windDirection, colornames.Red)
		LayLine(rc.course, rc.MarkX, rc.MarkY, TackAngle+rc.windDirection, colornames.Green)
	}

	if rc.showWindDirection {
		LayLine(rc.course, rc.MarkX, rc.MarkY, rc.windDirection, colornames.Blueviolet)
	}

	DrawWindDirection(rc.course, 1024-50, 768-50, rc.windDirection)

	return rc.course
}
