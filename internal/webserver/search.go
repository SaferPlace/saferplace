package webserver

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"safer.place/internal/address"
)

// SearchResponse is returned when the user searches for something
type SearchResponse struct {
	Response

	Address                string
	CoordX, CoordY         float64
	RoundedScore           int
	Score                  float64
	DistanceToUniversities map[string]int
}

func (s *WebServer) search(c *gin.Context) {
	req, err := s.request(c)
	if err != nil {
		var serr Error
		if errors.As(err, &serr) {
			c.HTML(serr.Code, "error.html", serr)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	addr, x, y, err := s.addressResolver.Resolve(req.Query)
	// We want to skip the empty Unresolved, as this will mean the address
	// is empty and it will be handed by the template
	if err != nil && !errors.Is(err, address.ErrUnresolved) {
		c.HTML(http.StatusInternalServerError, "error.html", Error{
			Code:  http.StatusInternalServerError,
			Cause: err,
		})
		return
	}

	score := s.scorer.Score(x, y)

	c.HTML(http.StatusOK, "search.html", SearchResponse{
		Response: s.response(req),

		Address:      addr,
		CoordX:       x,
		CoordY:       y,
		Score:        score,
		RoundedScore: int(math.Round(score)),
	})
}
