package webserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *WebServer) about(c *gin.Context) {
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

	c.HTML(http.StatusOK, "about.html", s.response(req))

}
