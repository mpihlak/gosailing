package gosailing

import (
	"github.com/gopxl/pixel/v2"
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

// markDirection is the heading to the mark from the start line (ie. 0 is straight north)
// windDirection is the heading from where the wind is blowing from (ie. 0 means wind is from the north)
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
	// Port layline
	rc.course.Clear()

	// Flag of the mark
	rc.course.Color = colornames.Orangered
	rc.course.Push(pixel.V(rc.MarkX, rc.MarkY), pixel.V(rc.MarkX, rc.MarkY+10))
	rc.course.Line(2)
	rc.course.Push(pixel.V(rc.MarkX, rc.MarkY+10), pixel.V(rc.MarkX, rc.MarkY+20))
	rc.course.Push(pixel.V(rc.MarkX, rc.MarkY+10), pixel.V(rc.MarkX+10, rc.MarkY+15))
	rc.course.Push(pixel.V(rc.MarkX, rc.MarkY+20), pixel.V(rc.MarkX+10, rc.MarkY+15))
	rc.course.Line(2)
	rc.course.Push(pixel.V(rc.MarkX, rc.MarkY))
	rc.course.Circle(2, 2)

	if rc.laylines {
		// Port layline
		rc.course.Color = colornames.Red
		portX, portY := RotatePoint(rc.MarkX, 0, rc.MarkX, rc.MarkY, -TackAngle+rc.windDirection)
		rc.course.Push(pixel.V(rc.MarkX, rc.MarkY), pixel.V(portX, portY))
		rc.course.Line(2)
		// Starboard layline
		rc.course.Color = colornames.Green
		starboardX, starboardY := RotatePoint(rc.MarkX, 0, rc.MarkX, rc.MarkY, TackAngle+rc.windDirection)
		rc.course.Push(pixel.V(rc.MarkX, rc.MarkY), pixel.V(starboardX, starboardY))
		rc.course.Line(2)
	}

	if rc.showWindDirection {
		// Wind direction indicator
		wdStartX := rc.MarkX
		wdStartY := rc.MarkY
		wdEndX, wdEndY := RotatePoint(wdStartX, wdStartY-1000, wdStartX, wdStartY, rc.windDirection)
		rc.course.Color = colornames.Blueviolet
		rc.course.Push(pixel.V(wdStartX, wdStartY), pixel.V(wdEndX, wdEndY))
		rc.course.Line(1)
	}

	DrawWindDirection(rc.course, 1024-50, 768-50, -rc.windDirection)

	return rc.course
}
