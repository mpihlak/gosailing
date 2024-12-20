package gosailing

import (
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

// DrawFlag draws a flag for the mark
func DrawFlag(canvas *imdraw.IMDraw, x, y float64) {
	canvas.Color = colornames.Orangered
	canvas.Push(pixel.V(x, y), pixel.V(x, y+10))
	canvas.Line(2)
	canvas.Push(pixel.V(x, y+10), pixel.V(x, y+20))
	canvas.Push(pixel.V(x, y+10), pixel.V(x+10, y+15))
	canvas.Push(pixel.V(x, y+20), pixel.V(x+10, y+15))
	canvas.Line(2)
	canvas.Push(pixel.V(x, y))
	canvas.Circle(2, 2)
}

// DrawBoat draws a little triangle for the boat
func DrawBoat(canvas *imdraw.IMDraw, x, y, heading float64) {
	canvas.Color = colornames.Darkblue
	// bow
	bowX, bowY := RotatePoint(x, y+7.5, x, y, heading)
	canvas.Push(pixel.V(bowX, bowY))
	// aft starboard corner
	sbX, sbY := RotatePoint(x+5, y-7.5, x, y, heading)
	canvas.Push(pixel.V(sbX, sbY))
	// aft port corner
	pX, pY := RotatePoint(x-5, y-7.5, x, y, heading)
	canvas.Push(pixel.V(pX, pY))
	// back to bow
	canvas.Push(pixel.V(bowX, bowY))
	canvas.Polygon(2)
}

// LayLine draws a line from the specified point to the specified heading
func LayLine(canvas *imdraw.IMDraw, x, y, heading float64, color color.RGBA) {
	canvas.Color = color
	laylineX, laylineY := RotatePoint(x, y-1000, x, y, heading)
	canvas.Push(pixel.V(x, y), pixel.V(laylineX, laylineY))
	canvas.Line(2)
}

// DrawWindDirection draws a wind direction indicator - a circle with an arrow inside it
func DrawWindDirection(canvas *imdraw.IMDraw, x, y, heading float64) {
	windDirection := heading

	canvas.Color = colornames.Gray
	canvas.Push(pixel.V(x, y))
	canvas.Circle(30, 2)

	lineX := x
	lineTopY := y + 25
	lineBottomY := y - 25

	// Wind direction line
	startX, startY := RotatePoint(lineX, lineTopY, x, y, windDirection)
	endX, endY := RotatePoint(lineX, lineBottomY, x, y, windDirection)
	canvas.Push(pixel.V(startX, startY), pixel.V(endX, endY))

	// Arrowhead
	arrowTipX := lineX
	arrowTipY := lineBottomY
	arrowLeftY := lineBottomY + 5
	arrowLeftX := lineX - 5
	arrowRightX := lineX + 5
	arrowRightY := lineBottomY + 5
	canvas.Push(pixel.V(RotatePoint(arrowLeftX, arrowLeftY, x, y, windDirection)))
	canvas.Push(pixel.V(RotatePoint(arrowTipX, arrowTipY, x, y, windDirection)))
	canvas.Push(pixel.V(RotatePoint(arrowRightX, arrowRightY, x, y, windDirection)))

	canvas.Line(2)
}
