package gosailing

import (
	"math"
)

type WindShifter struct {
	baseDirection float64
	amplitude     float64
	period        float64
	clock         float64
}

func NewWindShifter(baseDirection, amplitude, period float64) *WindShifter {
	return &WindShifter{
		baseDirection: baseDirection,
		amplitude:     amplitude,
		period:        period,
	}
}

func (ws *WindShifter) GetWindDirection() float64 {
	shift := ws.amplitude * math.Sin(2*math.Pi*ws.clock/ws.period)
	ws.clock += 0.05
	return ws.baseDirection + shift
}
