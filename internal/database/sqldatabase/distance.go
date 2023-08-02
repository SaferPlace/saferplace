package sqldatabase

import "math"

const earthRadius = 6371009

// distance calculates the distance between two points on a globe.
// Adapted from https://en.wikipedia.org/wiki/Great-circle_distance
// Because we are operating on float64, we do not care about inprecision errors
// as we care about very small distances.
func distance(x1, y1, x2, y2 float64) float64 {
	// convert to radians
	x1, y1, x2, y2 = dtor(x1), dtor(y1), dtor(x2), dtor(y2)

	lonDiff := math.Abs(x2 - x1)

	a := math.Sin(y1) * math.Sin(y2)
	b := math.Cos(y1) * math.Cos(y2) * math.Cos(lonDiff)
	rd := math.Acos(a + b)

	return rd * earthRadius
}

// dtor converts degrees to radians
func dtor(d float64) float64 {
	return (d * math.Pi) / 180
}
