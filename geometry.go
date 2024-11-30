package gosailing

import "math"

const (
	TackAngle = 45.0
)

// RotatePoint rotates a point (x, y) by n degrees around the specified origin (ox, oy)
func RotatePoint(x, y, ox, oy, degrees float64) (float64, float64) {
	// Translate point back to origin
	x -= ox
	y -= oy

	// Convert degrees to radians
	radians := degrees * (math.Pi / 180.0)

	// Apply rotation matrix
	cosTheta := math.Cos(radians)
	sinTheta := math.Sin(radians)

	xNew := x*cosTheta - y*sinTheta
	yNew := x*sinTheta + y*cosTheta

	// Translate point back
	xNew += ox
	yNew += oy

	return xNew, yNew
}
