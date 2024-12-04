package gosailing

import (
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

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

func DrawBoat(canvas *imdraw.IMDraw, x, y, heading float64) {
	// Draw a little triangle for the boat
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

func LayLine(canvas *imdraw.IMDraw, x, y, heading float64, color color.RGBA) {
	canvas.Color = color
	starboardX, starboardY := RotatePoint(x, 0, x, y, heading)
	canvas.Push(pixel.V(x, y), pixel.V(starboardX, starboardY))
	canvas.Line(2)
}