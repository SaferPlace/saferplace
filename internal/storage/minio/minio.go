// Copyright 2023 SaferPlace

// Package minio allows to upload images directly to minio.
package minio

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"safer.place/internal/config/secret"
)

type Config struct {
	Endpoint  string        `yaml:"endpoint" default:"minio.svc"`
	Bucket    string        `yaml:"bucket" default:"images"`
	AccessKey string        `yaml:"access_key" split_words:"true"`
	SecretKey secret.Secret `yaml:"secret_key" split_words:"true"`
	Secure    bool          `yaml:"secure" default:"false"`
}

type Storage struct {
	client *minio.Client
	bucket string
	tracer trace.Tracer
}

func New(ctx context.Context, cfg *Config, opts ...Option) (*Storage, error) {
	var err error

	s := &Storage{bucket: cfg.Bucket}

	for _, opt := range opts {
		opt(s)
	}

	s.client, err = minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, string(cfg.SecretKey), ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create minio client: %w", err)
	}

	// Create the bucket if it doesn't exist
	exists, err := s.client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("unable to check does the bucket exist: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("unable to create bucket: %w", err)
		}
	}

	if err := validate(s); err != nil {
		return nil, fmt.Errorf("minio validation failed: %w", err)
	}

	return s, nil
}

// Upload image to the minio bucket
func (s *Storage) Upload(ctx context.Context, r io.Reader, size int64, contentType string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "upload")
	defer span.End()

	id := uuid.New().String()
	if _, err := s.client.PutObject(
		ctx, s.bucket, id, r, size, minio.PutObjectOptions{
			ContentType: contentType,
		},
	); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("unable to upload image: %w", err)
	}

	return id, nil
}

var (
	errMissingClient = errors.New("missing client")
	errMissingBucket = errors.New("missing bucket")
	errMissingTracer = errors.New("missing tracer")
)

func validate(s *Storage) error {
	if s.bucket == "" {
		return errMissingBucket
	}
	if s.client == nil {
		return errMissingClient
	}
	if s.tracer == nil {
		return errMissingTracer
	}
	return nil
}
