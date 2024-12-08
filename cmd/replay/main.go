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

	var replayPoints []datasource.NavigationDataPoint
	for {
		d, ok := replayData.Next()
		if !ok {
			break
		}
		replayPoints = append(replayPoints, d)
	}

	// Calculate lat/lon bounds from replay points
	minLat, maxLat := replayPoints[0].Latitude, replayPoints[0].Latitude
	minLng, maxLng := replayPoints[0].Longitude, replayPoints[0].Longitude

	for _, p := range replayPoints {
		if p.Latitude < minLat {
			minLat = p.Latitude
		}
		if p.Latitude > maxLat {
			maxLat = p.Latitude
		}
		if p.Longitude < minLng {
			minLng = p.Longitude
		}
		if p.Longitude > maxLng {
			maxLng = p.Longitude
		}
	}

	zoom := 14000.0
	minX, minY := gosailing.LatLngToScreen(minLat, minLng, zoom)
	maxX, maxY := gosailing.LatLngToScreen(maxLat, maxLng, zoom)
	xOffset := minX - 50
	yOffset := minY - 50

	finished := false
	paused := false
	currentIndex := 0
	canvas := imdraw.New(nil)
	backgroundCanvas := imdraw.New(nil)

	// Draw a rectangle for the bounds
	backgroundCanvas.Color = colornames.Gray
	backgroundCanvas.Push(
		pixel.V(minX, minY),
		pixel.V(maxX, maxY),
	)
	backgroundCanvas.Rectangle(2)
	if *markLat != 0 && *markLng != 0 {
		x, y := gosailing.LatLngToScreen(*markLat, *markLng, zoom)
		gosailing.DrawFlag(backgroundCanvas, x-xOffset, y-yOffset)
	}

	for !win.Closed() {
		if keyPressed(pixel.KeyQ) || keyPressed(pixel.KeyEscape) {
			break
		}
		if keyPressed(pixel.KeySpace) || keyPressed(pixel.KeyP) {
			paused = !paused
		}
		if keyPressed(pixel.KeyR) {
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
				currentIndex++

				canvas.Clear()
				x, y := gosailing.LatLngToScreen(d.Latitude, d.Longitude, zoom)
				fmt.Printf("heading: %f\n", d.Heading)
				gosailing.DrawBoat(canvas, x-xOffset, y-yOffset, -d.CourseOverGround)
				gosailing.DrawWindDirection(canvas, 1024-50, 768-50, d.TrueWindDirection)

				backgroundCanvas.Color = colornames.Blue
				backgroundCanvas.Push(pixel.V(x-xOffset, y-yOffset))
				backgroundCanvas.Circle(2, 1)
			}
		}

		win.Clear(colornames.Lightblue)
		backgroundCanvas.Draw(win)
		canvas.Draw(win)
		win.Update()
	}
}

func main() {
	flag.Parse()
	opengl.Run(run)
}
