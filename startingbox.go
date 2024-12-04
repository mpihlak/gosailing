package gosailing

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type StartingBox struct {
	BoatEndX          float64
	BoatEndY          float64
	PinEndX           float64
	PinEndY           float64
	windDirection     float64
	laylines          bool
	showWindDirection bool
	canvas            *imdraw.IMDraw
}

func NewStartingBox(boatX, boatY, pinX, pinY, windDirection float64) *StartingBox {
	return &StartingBox{
		BoatEndX:          boatX,
		BoatEndY:          boatY,
		PinEndX:           pinX,
		PinEndY:           pinY,
		windDirection:     windDirection,
		laylines:          true,
		showWindDirection: true,
		canvas:            imdraw.New(nil),
	}
}

// SetWindDirection sets the current wind direction in degrees
// North is 0 and is straight up.
func (rc *StartingBox) SetWindDirection(direction float64) {
	rc.windDirection = direction
}

func (rc *StartingBox) ToggleLaylines() {
	rc.laylines = !rc.laylines
}

func (rc *StartingBox) ToggleWindDirection() {
	rc.showWindDirection = !rc.showWindDirection
}

func (rc *StartingBox) Drawable() *imdraw.IMDraw {
	rc.canvas.Clear()

	DrawFlag(rc.canvas, rc.PinEndX, rc.PinEndY)
	DrawBoat(rc.canvas, rc.BoatEndX, rc.BoatEndY, 0)

	if rc.laylines {
		// Boat end laylines
		LayLine(rc.canvas, rc.BoatEndX, rc.BoatEndY, -TackAngle+rc.windDirection, colornames.Red)
		LayLine(rc.canvas, rc.BoatEndX, rc.BoatEndY, TackAngle+rc.windDirection, colornames.Green)

		// Pin end laylines
		LayLine(rc.canvas, rc.PinEndX, rc.PinEndY, -TackAngle+rc.windDirection, colornames.Red)
		LayLine(rc.canvas, rc.PinEndX, rc.PinEndY, TackAngle+rc.windDirection, colornames.Green)

		// Starting line
		rc.canvas.Color = colornames.Blue
		rc.canvas.Push(pixel.V(rc.PinEndX, rc.PinEndY), pixel.V(rc.BoatEndX, rc.BoatEndY))
		rc.canvas.Line(2)
	}

	if rc.showWindDirection {
		// Wind direction indicator
		wdStartX := (rc.PinEndX + rc.BoatEndX) / 2
		wdStartY := rc.PinEndX + 1000
		wdEndX, wdEndY := RotatePoint(wdStartX, wdStartY-1000, wdStartX, wdStartY, rc.windDirection)
		rc.canvas.Color = colornames.Blueviolet
		rc.canvas.Push(pixel.V(wdStartX, wdStartY), pixel.V(wdEndX, wdEndY))
		rc.canvas.Line(1)
	}

	return rc.canvas
}
