package score

// Scorer creates a safety score based on a the given location.
type Scorer interface {
	Score(x, y float64) float64
}
