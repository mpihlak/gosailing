package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.io/mpihlak/gosailing"
	"github.io/mpihlak/gosailing/datasource"
	"golang.org/x/image/colornames"
)

const (
	maxWidth  = 1024.0
	maxHeight = 768.0

	boatLocationX = maxWidth / 2
	boatLocationY = 25
)

var (
	csvFile   = flag.String("csv", "", "CSV data file to replay")
	startTime = flag.String("start", "", "Start time to replay from (RFC3339 format)")
	endTime   = flag.String("end", "", "End time to replay to (RFC3339 format)")
	markLat   = flag.Float64("markLat", 0, "Latitude of the mark")
	markLng   = flag.Float64("markLng", 0, "Longitude of the mark")
)

func run() {
	if *csvFile == "" {
		log.Fatalf("Must provide -csv argument with replay file")
	}
	if *markLat == 0 || *markLng == 0 {
		log.Fatalf("Must provide -markLat and -markLng arguments with mark location")
	}

	var start, end time.Time
	var err error
	if *startTime != "" {
		start, err = time.Parse(time.RFC3339, *startTime)
		if err != nil {
			log.Fatalf("Invalid start time format: %v", err)
		}
	}
	if *endTime != "" {
		end, err = time.Parse(time.RFC3339, *endTime)
		if err != nil {
			log.Fatalf("Invalid end time format: %v", err)
		}
	}

	f, err := os.Open(*csvFile)
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}

	replayData, err := datasource.NewReplayNavigationDataProvider(f, &start, &end)
	if err != nil {
		log.Fatal("Unable to load replay")
	}

	cfg := opengl.WindowConfig{
		Title:  "Go Sailing!",
		Bounds: pixel.R(0, 0, maxWidth, maxHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Throttle the keyboard to avoid registering unintended repeated keypresses
	lastKeyPressed := make(map[pixel.Button]time.Time)
	keyPressed := func(k pixel.Button) bool {
		if win.Pressed(k) && time.Since(lastKeyPressed[k]) > 200*time.Millisecond {
			lastKeyPressed[k] = time.Now()
			return true
		}
		return false
	}

	replayPoints := replayData.GetAllPoints()

	medianWind := datasource.MedianWindDirection(replayPoints)
	fmt.Printf("Median wind direction: %.2f\n", medianWind)

	zoom := 5500.0

	markX, markY := gosailing.LatLngToScreen(*markLat, *markLng, zoom)

	type locationXY struct {
		x float64
		y float64
	}

	locations := make([]locationXY, len(replayPoints))
	var minX, minY, maxX, maxY float64
	for i, p := range replayPoints {
		x, y := gosailing.LatLngToScreen(p.Latitude, p.Longitude, zoom)
		if *markLat != 0 && *markLng != 0 {
			x, y = gosailing.RotatePoint(x, y, markX, markY, -medianWind)
		}
		replayPoints[i].CourseOverGround -= medianWind
		replayPoints[i].TrueWindDirection -= medianWind
		locations[i] = locationXY{x, y}
		if i == 0 || x < minX {
			minX = x
		}
		if i == 0 || y < minY {
			minY = y
		}
		if i == 0 || x > maxX {
			maxX = x
		}
		if i == 0 || y > maxY {
			maxY = y
		}
	}

	fmt.Printf("minX: %.2f, minY: %.2f, maxX: %.2f, maxY: %.2f\n", minX, minY, maxX, maxY)
	xOffset := markX - maxWidth/2
	yOffset := minY - 50

	finished := false
	paused := false
	currentIndex := 0
	canvas := imdraw.New(nil)
	trackCanvas := imdraw.New(nil)

	for !win.Closed() {
		if keyPressed(pixel.KeyQ) || keyPressed(pixel.KeyEscape) {
			break
		}
		if keyPressed(pixel.KeySpace) || keyPressed(pixel.KeyP) {
			paused = !paused
		}
		if keyPressed(pixel.KeyR) {
			trackCanvas.Clear()
			currentIndex = 0
			finished = false
		}
		if paused {
			time.Sleep(100 * time.Millisecond)
		} else if !finished {
			if currentIndex >= len(replayPoints) {
				log.Printf("Replay data source finished")
				finished = true
			} else {
				d := replayPoints[currentIndex]
				l := locations[currentIndex]
				currentIndex++

				canvas.Clear()

				gosailing.DrawBoat(canvas, l.x-xOffset, l.y-yOffset, d.CourseOverGround)
				gosailing.LayLine(canvas, l.x-xOffset, l.y-yOffset, d.TrueWindDirection+45+180, colornames.Red)
				gosailing.LayLine(canvas, l.x-xOffset, l.y-yOffset, d.TrueWindDirection-45+180, colornames.Green)
				gosailing.LayLine(canvas, l.x-xOffset, l.y-yOffset, d.CourseOverGround+180, colornames.Gray)

				gosailing.DrawWindDirection(canvas, 1024-50, 768-50, d.TrueWindDirection)

				// draw track
				trackCanvas.Color = colornames.Blueviolet
				trackCanvas.Push(pixel.V(l.x-xOffset, l.y-yOffset))
				trackCanvas.Circle(1, 1)

				gosailing.DrawFlag(canvas, markX-xOffset, markY-yOffset)
				gosailing.LayLine(canvas, markX-xOffset, markY-yOffset, d.TrueWindDirection+45, colornames.Red)
				gosailing.LayLine(canvas, markX-xOffset, markY-yOffset, d.TrueWindDirection-45, colornames.Green)

				time.Sleep(50 * time.Millisecond)
			}
		}

		win.Clear(colornames.Lightblue)
		canvas.Draw(win)
		trackCanvas.Draw(win)
		win.Update()
	}
}

func main() {
	flag.Parse()
	opengl.Run(run)
}
