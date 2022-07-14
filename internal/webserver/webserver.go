package webserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"safer.place/internal/address"
	"safer.place/internal/language"
	"safer.place/internal/score"
)

// WebServer serves the website
type WebServer struct {
	router *gin.Engine

	addressResolver address.Resolver
	scorer          score.Scorer

	// prepared language information
	languages       map[string]language.Language
	languageOptions []language.Info

	// Fonts
	NormalFont, FancyFont string
}

func New(opts ...Option) *WebServer {
	s := &WebServer{
		router:    gin.Default(),
		languages: make(map[string]language.Language),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.router.GET("/search", s.search)
	s.router.GET("/about", s.about)
	s.router.GET("/", s.index)

	return s
}

func (s *WebServer) Run(port int) error {
	return s.router.Run(fmt.Sprintf(":%d", port))
}

// Request contains the common fields which might be present in every request.
// Individual requests might embed this type.
type Request struct {
	Language string `form:"lang"`
	Query    string `form:"q"` // We allow the query on every page.
}

// Response contains the common fields present in every response. Individual
// responses typically embed this type.
type Response struct {
	// CriticalError that can be returned if something goes critically wrong.
	// This means that we might not be able to even return language specific
	// response.
	CriticalError string

	Query string

	// AvailableLanguages that can be used
	AvailableLanguages []language.Info
	// Lang that that is currently selected
	Lang string
	// Language data which contains the actual messages.
	language.Language

	NormalFont string
	FancyFont  string
}

// reques parses the standard request from the context.
func (s *WebServer) request(c *gin.Context) (Request, error) {
	var req Request
	if err := c.ShouldBind(&req); err != nil {
		return req, fmt.Errorf("cannot bind: %w", err)
	}

	// default language to english
	if _, ok := s.languages[req.Language]; !ok {
		req.Language = "en"
	}

	return req, nil
}

// response builds the common Response data from the context and the request.
func (s *WebServer) response(req Request) Response {
	return Response{
		NormalFont: s.NormalFont,
		FancyFont:  s.FancyFont,

		Query: req.Query,

		AvailableLanguages: s.languageOptions,
		// This places trust that the language in the request is valid.
		// Do we want this?
		Lang:     req.Language,
		Language: s.languages[req.Language],
	}
}
