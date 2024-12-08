package main

import (
	"flag"
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

	finished := false
	canvas := imdraw.New(nil)
	for !win.Closed() {
		if keyPressed(pixel.KeyQ) || keyPressed(pixel.KeyEscape) {
			break
		}

		if !finished {
			d, ok := replayData.Next()

			if !ok {
				log.Printf("Replay data source finished")
				finished = true
			} else {
				canvas.Clear()
				gosailing.DrawWindDirection(canvas, 1024-50, 768-50, d.TrueWindDirection)
				log.Printf("data point: %v", d)
			}
		}

		win.Clear(colornames.Lightblue)
		canvas.Draw(win)
		win.Update()
	}
}

func main() {
	flag.Parse()
	opengl.Run(run)
}
