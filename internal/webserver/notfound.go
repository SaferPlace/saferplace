package webserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundResponse is returned for all 404 pages.
type NotFoundResponse struct {
	Response

	ErrorMessage string
}

func (s *WebServer) notfound(c *gin.Context) {
	status := http.StatusNotFound
	msg := "404 Not Found"
	if c.Query("go-get") == "1" {
		status = http.StatusOK
		msg = "Let's Go!"
	}
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

	c.HTML(status, "notfound.html", NotFoundResponse{
		ErrorMessage: msg,
		Response:     s.response(req),
	})
}
