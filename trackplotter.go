package gosailing

import (
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type TrackPlotter struct {
	track        *imdraw.IMDraw
	lastPlotTime time.Time
}

func NewTrackPlotter() *TrackPlotter {
	return &TrackPlotter{
		track: imdraw.New(nil),
	}
}

func (tp *TrackPlotter) PlotLocation(x, y float64) {
	if time.Since(tp.lastPlotTime) > 100*time.Millisecond {
		tp.track.Color = colornames.Blueviolet
		tp.track.Push(pixel.V(x, y))
		tp.track.Circle(2, 1)
		tp.lastPlotTime = time.Now()
	}
}

func (tp *TrackPlotter) Drawable() *imdraw.IMDraw {
	return tp.track
}
