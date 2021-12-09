package stations

import (
	"log"
	"math"
	"sort"

	"safer.place/data"
)

type Stations struct {
	stations map[string]*Station
}

func New() *Stations {
	ss := &Stations{
		stations: make(map[string]*Station),
	}

	for name, coords := range data.GardaLocations() {
		ss.stations[name] = newStation(name, coords.X, coords.Y)
	}

	// TODO: Add crimes
	for name, crimes := range data.CrimesPerStation() {
		// If we don't find a station with the given name, we create a new one,
		// but we provide coordinates as 0,0 since we have no idea where it is
		s, exists := ss.stations[name]
		if !exists {
			log.Printf("no coordinates for %q, skipping", name)
			continue
		}
		s.crimes = crimes
		ss.stations[name] = s
	}

	return ss
}

// Nearest n stations from the given coordinates
func (s *Stations) Nearest(x, y float64, n int) []*Station {
	bd := ByDistance{
		stations: make([]*Station, 0, len(s.stations)),
		x:        x,
		y:        y,
	}
	for _, s := range s.stations {
		bd.stations = append(bd.stations, s)
	}
	sort.Sort(bd)

	return bd.stations[:n]
}

// SafestAndDangerousScores gets the safest and most most dangerous score
func (s *Stations) SafestAndDangerousScores(years int) (float64, float64) {
	min, max := math.Inf(1), 0.0
	minName, maxName := "", ""
	for _, s := range s.stations {
		score := s.ScoreAverage(years)
		// Somehow ignore the 0.0, no idea why its showing up.
		if score < 0.01 {
			continue
		}
		if score < min {
			min = score
			minName = s.Name
			continue
		}
		if max < score {
			max = score
			maxName = s.Name
			continue
		}
	}

	log.Printf("Safest: %q (%.2f) - Most Dangerous: %q (%.2f)",
		minName, min, maxName, max)

	return min, max
}

type ByDistance struct {
	x, y     float64
	stations []*Station
}

func (ss ByDistance) Len() int { return len(ss.stations) }
func (ss ByDistance) Less(i, j int) bool {
	return ss.stations[i].DistanceTo(ss.x, ss.y) <
		ss.stations[j].DistanceTo(ss.x, ss.y)
}
func (ss ByDistance) Swap(i, j int) {
	ss.stations[i], ss.stations[j] = ss.stations[j], ss.stations[i]
}