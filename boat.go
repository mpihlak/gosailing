package gosailing

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type Boat struct {
	currentX      float64
	currentY      float64
	heading       float64
	windDirection float64
	boat          *imdraw.IMDraw
}

func NewBoat(currentX, currentY, windDirection float64) *Boat {
	// Initial heading is starboard tack, going upwind
	heading := windDirection - TackAngle

	return &Boat{
		currentX:      currentX,
		currentY:      currentY,
		heading:       heading,
		windDirection: windDirection,
		boat:          imdraw.New(nil),
	}
}

func (b *Boat) SetWindDirection(direction float64) {
	b.windDirection = direction

	// Drive as close to the wind as possible
	if b.heading < b.windDirection {
		b.heading = b.windDirection - TackAngle
	} else {
		b.heading = b.windDirection + TackAngle
	}
}

func (b *Boat) Tack() {
	if b.heading < b.windDirection {
		b.heading = b.windDirection + TackAngle
	} else {
		b.heading = b.windDirection - TackAngle
	}
}

func (b *Boat) Advance() {
	newX, newY := RotatePoint(b.currentX, b.currentY+1, b.currentX, b.currentY, b.heading)
	b.currentX = newX
	b.currentY = newY
}

func (b *Boat) GetXY() (float64, float64) {
	return b.currentX, b.currentY
}

func (b *Boat) Drawable() *imdraw.IMDraw {
	b.boat.Clear()
	b.boat.Color = colornames.Darkblue
	b.boat.Push(pixel.V(b.currentX, b.currentY))
	b.boat.Circle(10.0, 5)

	// Starboard layline
	b.boat.Color = colornames.Green
	sbLayLineX, sbLayLineY := RotatePoint(b.currentX, b.currentY+1000, b.currentX, b.currentY, b.windDirection-TackAngle)
	b.boat.Push(pixel.V(b.currentX, b.currentY), pixel.V(sbLayLineX, sbLayLineY))
	b.boat.Line(2)

	// Port layline
	b.boat.Color = colornames.Red
	portLayLineX, portLayLineY := RotatePoint(b.currentX, b.currentY+1000, b.currentX, b.currentY, b.windDirection+TackAngle)
	b.boat.Push(pixel.V(b.currentX, b.currentY), pixel.V(portLayLineX, portLayLineY))
	b.boat.Line(2)

	return b.boat
}
