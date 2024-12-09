package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.io/mpihlak/gosailing"
	"github.io/mpihlak/gosailing/datasource"
	"golang.org/x/image/colornames"
)

const (
	maxWidth  = 1024.0
	maxHeight = 768.0
)

var (
	csvFile   = flag.String("csv", "", "CSV data file to replay")
	startTime = flag.String("start", "", "Start time to replay from (RFC3339 format)")
	endTime   = flag.String("end", "", "End time to replay to (RFC3339 format)")
	markLat   = flag.Float64("markLat", 0, "Latitude of the mark")
	markLng   = flag.Float64("markLng", 0, "Longitude of the mark")
	zoomLevel = flag.Float64("zoom", 5500, "Zoom level")
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

	rr, err := gosailing.NewRaceReplay(*markLat, *markLng, maxWidth, maxHeight, *zoomLevel, replayData)
	if err != nil {
		log.Fatalf("Unable to create race replay: %v", err)
	}

	for !win.Closed() {
		if keyPressed(pixel.KeyQ) || keyPressed(pixel.KeyEscape) {
			break
		}
		if keyPressed(pixel.KeySpace) || keyPressed(pixel.KeyP) {
			rr.TogglePause()
		}
		if keyPressed(pixel.KeyR) {
			rr.StartReplay()
		}

		win.Clear(colornames.Lightblue)
		rr.Update(win)
		win.Update()
	}
}

func main() {
	flag.Parse()
	opengl.Run(run)
}
