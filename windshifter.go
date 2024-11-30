package gosailing

import (
	"math"
	"time"
)

type WindShifter struct {
	baseDirection float64
	amplitude     float64
	period        float64
	startTime     time.Time
}

func NewWindShifter(baseDirection, amplitude, period float64) *WindShifter {
	return &WindShifter{
		baseDirection: baseDirection,
		amplitude:     amplitude,
		period:        period,
		startTime:     time.Now(),
	}
}

func (ws *WindShifter) GetWindDirection() float64 {
	elapsed := time.Since(ws.startTime).Seconds()
	shift := ws.amplitude * math.Sin(2*math.Pi*elapsed/ws.period)
	return ws.baseDirection + shift
}
