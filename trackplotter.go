package gosailing

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type TrackPlotter struct {
	track    *imdraw.IMDraw
	plottedX float64
	plottedY float64
}

func NewTrackPlotter(x, y float64) *TrackPlotter {
	return &TrackPlotter{
		plottedX: x,
		plottedY: y,
		track:    imdraw.New(nil),
	}
}

func (tp *TrackPlotter) PlotLocation(x, y float64) {
	distance := math.Hypot(tp.plottedX-x, tp.plottedY-y)
	if distance > 5 {
		tp.track.Color = colornames.Blueviolet
		tp.track.Push(pixel.V(tp.plottedX, tp.plottedY))
		tp.track.Circle(1, 1)
		tp.plottedX = x
		tp.plottedY = y
	}
}

func (tp *TrackPlotter) Drawable() *imdraw.IMDraw {
	return tp.track
}
