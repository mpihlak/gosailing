package gosailing

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

type WindShifter interface {
	GetWindDirection() float64
}

type ReplayWindShifter struct {
	windDirections []float64
	pos            float64
}

func NewReplayShifter(fileName string) *ReplayWindShifter {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	windDirections := make([]float64, 0)
	for _, line := range strings.Split(string(buf), "\n") {
		if line == "" {
			continue
		}

		wd, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Printf("ignoring invalid wind direction: %v\n", line)
			continue
		}

		windDirections = append(windDirections, wd)
	}

	if len(windDirections) < 1 {
		panic("no wind data")
	}

	// Normalize wind directions to north

	sortedDirections := make([]float64, len(windDirections))
	copy(sortedDirections, windDirections)

	sort.Float64s(sortedDirections)
	n := len(sortedDirections)
	medianWind := sortedDirections[n/2]
	if len(sortedDirections)%2 == 0 {
		medianWind = (sortedDirections[n/2] + sortedDirections[n/2+1]) / 2
	}

	for i := range windDirections {
		windDirections[i] -= medianWind
	}

	return &ReplayWindShifter{
		windDirections: windDirections,
		pos:            rand.Float64() * float64(len(windDirections)),
	}
}

// TODO: Use the commonly understood wind direction, ie. where is it blowing from
func (r *ReplayWindShifter) GetWindDirection() float64 {
	val := r.windDirections[int(r.pos)]
	r.pos = r.pos + 0.05
	if int(r.pos) >= len(r.windDirections) {
		r.pos = 0
	}
	return val
}

type OscillatingWindShifter struct {
	baseDirection float64
	amplitude     float64
	period        float64
	shiftRate     float64
	clock         float64
}

func NewOscillatingWindShifter(baseDirection, amplitude, period, shiftRate float64) *OscillatingWindShifter {
	return &OscillatingWindShifter{
		baseDirection: baseDirection,
		amplitude:     amplitude,
		period:        period,
		shiftRate:     shiftRate,
	}
}

func (ws *OscillatingWindShifter) GetWindDirection() float64 {
	shift := ws.amplitude * math.Sin(2*math.Pi*ws.clock/ws.period)
	ws.clock += 0.05
	ws.baseDirection += ws.shiftRate
	return ws.baseDirection + shift
}
