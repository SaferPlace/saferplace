package webserver

import (
	"html/template"

	"safer.place/internal/address"
	"safer.place/internal/language"
	"safer.place/internal/score"
)

// Option modifies the behaviour of the webserver
type Option func(s *WebServer)

// Languages adds the language mapping to the webserver
func Languages(m map[language.Info]language.Language) Option {
	return func(s *WebServer) {
		for info, lang := range m {
			s.languages[info.Code] = lang
			s.languageOptions = append(s.languageOptions, info)
		}
	}
}

// Templates to be used for rendering
func Templates(tmpl *template.Template) Option {
	return func(s *WebServer) {
		s.router.SetHTMLTemplate(tmpl)
	}
}

// AddressResolver allows the webserver to resolve the address.
func AddressResolver(addrResolver address.Resolver) Option {
	return func(s *WebServer) {
		s.addressResolver = addrResolver
	}
}

// Scorer returns the score for the location
func Scorer(scorer score.Scorer) Option {
	return func(s *WebServer) {
		s.scorer = scorer
	}
}
