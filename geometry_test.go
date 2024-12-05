package gosailing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRotatePoint(t *testing.T) {
	require := require.New(t)

	// Clockwise 90 degrees
	rx, ry := RotatePoint(10, 0, 0, 0, 90)
	require.InDelta(float64(0.0), rx, 0.001)
	require.InDelta(float64(10.0), ry, 0.001)

	// Anti-clockwise 90 degrees
	rx, ry = RotatePoint(10, 0, 0, 0, -90)
	require.InDelta(float64(0.0), rx, 0.001)
	require.InDelta(float64(-10.0), ry, 0.001)
}
