package gosailing

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type TrackPlotter struct {
	canvas   *imdraw.IMDraw
	plottedX float64
	plottedY float64
}

func NewTrackPlotter(x, y float64) *TrackPlotter {
	return &TrackPlotter{
		plottedX: x,
		plottedY: y,
		canvas:   imdraw.New(nil),
	}
}

func (tp *TrackPlotter) PlotLocation(x, y float64) {
	distance := math.Hypot(tp.plottedX-x, tp.plottedY-y)
	if distance > 5 {
		tp.canvas.Color = colornames.Blueviolet
		tp.canvas.Push(pixel.V(tp.plottedX, tp.plottedY))
		tp.canvas.Circle(1, 1)
		tp.plottedX = x
		tp.plottedY = y
	}
}

func (tp *TrackPlotter) Clear() {
	tp.canvas.Clear()
}

func (tp *TrackPlotter) Drawable() *imdraw.IMDraw {
	return tp.canvas
}
