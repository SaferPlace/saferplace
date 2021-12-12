package v1

import "safer.place/internal/stations"

type Scorer struct {
	stations       *stations.Stations
	years, nearest int
	min, max       float64
}

func New(stations *stations.Stations) *Scorer {
	min, max := stations.SafestAndDangerousScores(5)
	return &Scorer{
		stations: stations,
		nearest:  3,
		years:    5,
		min:      min,
		max:      max,
	}
}

func (s *Scorer) Score(x, y float64) float64 {
	// TODO: Improve this algorithm. For now we get the safest location
	//	     and the worst location, and the score would fall somewhere
	//       in between. But we are definatelly calculating it wrong.
	nearest := s.stations.Nearest(x, y, s.nearest)

	sum := 0.0
	for _, station := range nearest {
		sum += station.ScoreAverage(s.years)
	}
	avg := sum / float64(s.nearest)

	// Score between 1-5
	score := 1 + ((avg - s.min) / (s.max - s.min) * 4)

	return score
}
