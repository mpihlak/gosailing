package gosailing

import (
	"math"

	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type Boat struct {
	currentX       float64
	currentY       float64
	heading        float64
	windDirection  float64
	sailedDistance float64
	laylines       bool
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
		laylines:      true,
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

func (b *Boat) ToggleLaylines() {
	b.laylines = !b.laylines
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

func (b *Boat) SetLocation(x, y, heading, windDirection float64) {
	b.sailedDistance += math.Hypot(b.currentX-x, b.currentY-y)
	b.currentX = x
	b.currentY = y
	b.heading = heading
	b.windDirection = windDirection
}

func (b *Boat) GetSailedDistance() float64 {
	return b.sailedDistance
}

func (b *Boat) Drawable() *imdraw.IMDraw {
	b.boat.Clear()

	DrawBoat(b.boat, b.currentX, b.currentY, b.heading)

	if b.laylines {
		LayLine(b.boat, b.currentX, b.currentY, b.windDirection+TackAngle+180, colornames.Red)
		LayLine(b.boat, b.currentX, b.currentY, b.windDirection-TackAngle+180, colornames.Green)
		LayLine(b.boat, b.currentX, b.currentY, b.heading+180, colornames.Gray)
	}

	return b.boat
}
