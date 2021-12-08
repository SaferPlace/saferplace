package stations

import (
	"math"

	"safer.place/data"
)

// ScoreMarker is a function which convers years of data into a score
type ScoreMarkerFunc func(years int) float64

type Station struct {
	Name   string
	X, Y   float64
	crimes data.CrimesPerType

	markerFns     []ScoreMarkerFunc
	currentYearFn func() int
}

func (s Station) String() string {
	return s.Name
}

func newStation(name string, x, y float64) *Station {
	s := &Station{
		Name: name,
		X:    x,
		Y:    y,
		currentYearFn: func() int {
			// Until we have more up to date data, pretend its 2016
			return 2016
		},
	}

	s.markerFns = []ScoreMarkerFunc{
		s.violentCrimeAverage,
		s.bulgaryAverage,
		s.theftAverage,
		s.robberyAverage,
	}

	return s
}

func (s *Station) violentCrimeAverage(years int) float64 {
	return s.avgCrimes(s.crimes[data.MurderAndAssaultCrime], years)
}

func (s *Station) bulgaryAverage(years int) float64 {
	return s.avgCrimes(s.crimes[data.BurglaryCrime], years)
}

func (s *Station) theftAverage(years int) float64 {
	return s.avgCrimes(s.crimes[data.TheftCrime], years)
}

func (s *Station) robberyAverage(years int) float64 {
	return s.avgCrimes(s.crimes[data.RobberyCrime], years)
}

func (s *Station) ScoreAverage(years int) float64 {
	sum := 0.0
	for _, marker := range s.markerFns {
		ca := marker(years)
		sum += ca
	}
	return sum / float64(len(s.markerFns))
}

// DistanceTo returns the approximate distance between coordinates
func (s *Station) DistanceTo(x, y float64) float64 {
	// approximate distance calculation
	dX, dY := s.X-x, s.Y-y
	return math.Sqrt(dX*dX + dY*dY)
}

func (s *Station) avgCrimes(crimes data.CrimesInYear, years int) float64 {
	from := s.currentYearFn() - years

	total := 0
	for i := from; i < s.currentYearFn(); i++ {
		total += crimes[i]
	}

	return float64(total) / float64(years)
}
