// Copyright 2023 SaferPlace

// Package imageupload allows for HTTP image uploads to a cloud storage.
package imageupload

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"safer.place/internal/storage"
)

// Service is the image upload service
type Service struct {
	storage storage.Storage
	log     *zap.Logger
}

// Register registers the image upload service.
func Register(log *zap.Logger, store storage.Storage) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return "/v1/upload", &Service{
			log:     log,
			storage: store,
		}
	}
}

// ServeHTTP is the handler accepting the image upload.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.log.Error("unable to parse form data", zap.Error(err))
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		s.log.Error("unable to get image from user", zap.Error(err))
		http.Error(w, "image upload failed", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reference, err := s.storage.Upload(
		r.Context(), file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		s.log.Error("image upload failed", zap.Error(err))
		http.Error(w, "image upload failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, reference)
}
