// Copyright 2023 SaferPlace

// Package imageupload allows for HTTP image uploads to a cloud storage.
package imageupload

import (
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"safer.place/internal/service"
	"safer.place/internal/storage"
)

// Service is the image upload service
type Service struct {
	tracer  trace.Tracer
	storage storage.Storage
	log     *zap.Logger
}

// Register registers the image upload service.
func Register(opts ...Option) service.Service {
	s := &Service{}

	for _, opt := range opts {
		opt(s)
	}

	if err := validate(s); err != nil {
		panic(err)
	}

	// We can ignore the interceptors as this is a non-connect service
	return func(_ ...connect.Interceptor) (string, http.Handler) {
		return "/v1/upload", s
	}
}

// ServeHTTP is the handler accepting the image upload.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.tracer.Start(r.Context(), "upload")
	defer span.End()
	if err := r.ParseForm(); err != nil {
		s.log.Error("unable to parse form data", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		s.log.Error("unable to get image from user", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "image upload failed", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reference, err := s.storage.Upload(
		ctx, file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		s.log.Error("image upload failed", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "image upload failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, reference)
}

var (
	errMissingLogger  = errors.New("missing logger")
	errMissingTrace   = errors.New("missing tracer")
	errMissingStorage = errors.New("missing storage")
)

func validate(s *Service) error {
	if s.log == nil {
		return errMissingLogger
	}
	if s.tracer == nil {
		return errMissingTrace
	}
	if s.storage == nil {
		return errMissingStorage
	}
	return nil
}
