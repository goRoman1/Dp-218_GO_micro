package models

import "math"

//Coordinate is for comfortable storage latitude, longitude values.
type Coordinate struct {
	Latitude 		float64 		`json:"latitude"`
	Longitude 		float64 		`json:"longitude"`
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// http://en.wikipedia.org/wiki/Haversine_formula
func (l1 Coordinate) Distance(l2 Coordinate) float64 {
	// convert to radians
	var la1, lo1, la2, lo2, earthRadius float64
	la1 = l1.Latitude * math.Pi / 180
	lo1 = l1.Longitude * math.Pi / 180
	la2 = l2.Latitude * math.Pi / 180
	lo2 = l2.Longitude * math.Pi / 180

	earthRadius = 6378100 // Earth radius in meters

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}