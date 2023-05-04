package oauth

import "net/http"

type Service struct {
	mux *http.ServeMux
}

func Register() (string, http.Handler) {
	s := &Service{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("github", s.githubCallback)

	return "/oauth", s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) githubCallback(w http.ResponseWriter, r *http.Request) {

}
