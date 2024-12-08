package gosailing

import "math"

const (
	TackAngle = 45.0
)

// RotatePoint rotates a point (x, y) by n degrees around the specified origin (ox, oy)
func RotatePoint(x, y, ox, oy, degrees float64) (float64, float64) {
	x -= ox
	y -= oy

	radians := -degrees * (math.Pi / 180.0)

	cosTheta := math.Cos(radians)
	sinTheta := math.Sin(radians)

	xNew := x*cosTheta - y*sinTheta
	yNew := x*sinTheta + y*cosTheta

	xNew += ox
	yNew += oy

	return xNew, yNew
}

// LatLngToScreen converts latitude and longitude to screen coordinates. The zoom value is approximately the
// number of pixels per degree of latitude and longitude. One nautical mile is 1/60 degrees of latitude,
// so for zoom=1, 60 nautical miles would be ~1 pixel (0.0166 pixels per nm). If we want 1 nm to be 100 pixels,
// we need zoom=100/0.0166=5993.6
func LatLngToScreen(latitude, longitude float64, zoom float64) (float64, float64) {
	const tileSize = 512

	sinY := math.Sin(toRadians(latitude))
	sinY = math.Min(math.Max(sinY, -0.9999), 0.9999)

	return (tileSize*(0.5+longitude/360) + 100) * zoom,
		(tileSize * (0.5 + math.Log((1+sinY)/(1-sinY))/(4*math.Pi))) * zoom
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
