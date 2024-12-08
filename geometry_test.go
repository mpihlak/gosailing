package gosailing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRotatePoint(t *testing.T) {
	require := require.New(t)

	// Clockwise 90 degrees
	rx, ry := RotatePoint(0, 10, 0, 0, 90)
	require.InDelta(float64(10.0), rx, 0.001)
	require.InDelta(float64(00.0), ry, 0.001)

	// Clockwise 180 degrees
	rx, ry = RotatePoint(10, 0, 0, 0, 90)
	require.InDelta(float64(0.0), rx, 0.001)
	require.InDelta(float64(-10.0), ry, 0.001)

	// Anti-clockwise 90 degrees
	rx, ry = RotatePoint(10, 0, 0, 0, 90)
	require.InDelta(float64(0.0), rx, 0.001)
	require.InDelta(float64(-10.0), ry, 0.001)
}

func TestLatLngToScreen(t *testing.T) {
	require := require.New(t)

	lat, lng := 59.0, 24.0
	zoom := 1.0
	x, y := LatLngToScreen(lat, lng, zoom)
	fmt.Printf("x: %f, y: %f\n", x, y)
	require.InDelta(float64(390.13), x, 0.01)
	require.InDelta(float64(360.51), y, 0.01)

	// Latitude changes by one degree, expect a ~1 pixel difference in Y
	x, y = LatLngToScreen(lat+1, lng, zoom)
	fmt.Printf("lat+1 x: %f, y: %f\n", x, y)
	require.InDelta(float64(390.13), x, 0.01)
	require.InDelta(float64(363.31), y, 0.01)

	// Longitude changes by one degree, expect a ~1 pixel difference in X
	x, y = LatLngToScreen(lat, lng+1, zoom)
	fmt.Printf("lng+1 x: %f, y: %f\n", x, y)
	require.InDelta(float64(391.55), x, 0.01)
	require.InDelta(float64(360.51), y, 0.01)
}
