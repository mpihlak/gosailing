package gosailing

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type Boat struct {
	currentX       float64
	currentY       float64
	heading        float64
	windDirection  float64
	sailedDistance float64
	boat           *imdraw.IMDraw
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
	b.sailedDistance += math.Hypot(b.currentX-newX, b.currentY-newY)
	b.currentX = newX
	b.currentY = newY
}

func (b *Boat) GetXY() (float64, float64) {
	return b.currentX, b.currentY
}

func (b *Boat) GetSailedDistance() float64 {
	return b.sailedDistance
}

func (b *Boat) Drawable() *imdraw.IMDraw {
	b.boat.Clear()

	// Draw a little triangle for the boat
	b.boat.Color = colornames.Darkblue
	// bow
	bowX, bowY := RotatePoint(b.currentX, b.currentY+7.5, b.currentX, b.currentY, b.heading)
	b.boat.Push(pixel.V(bowX, bowY))
	// aft starboard corner
	sbX, sbY := RotatePoint(b.currentX+5, b.currentY-7.5, b.currentX, b.currentY, b.heading)
	b.boat.Push(pixel.V(sbX, sbY))
	// aft port corner
	pX, pY := RotatePoint(b.currentX-5, b.currentY-7.5, b.currentX, b.currentY, b.heading)
	b.boat.Push(pixel.V(pX, pY))
	// back to bow
	b.boat.Push(pixel.V(bowX, bowY))
	b.boat.Polygon(2)

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
