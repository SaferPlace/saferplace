// Copyright 2023 SaferPlace

// Package minio allows to upload images directly to minio.
package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint  string `default:"minio.svc"`
	Bucket    string `default:"images"`
	AccessKey string `split_words:"true"`
	SecretKey string `split_words:"true"`
	Secure    bool   `default:"false"`
}

type Storage struct {
	client *minio.Client
	bucket string
}

func New() (*Storage, error) {
	var cfg Config
	if err := envconfig.Process("MINIO", &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse minio config: %w", err)
	}

	c, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create minio client: %w", err)
	}

	ctx := context.Background()
	// Create the bucket if it doesn't exist
	exists, err := c.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("unable to check does the bucket exist: %w", err)
	}
	if !exists {
		if err := c.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("unable to create bucket: %w", err)
		}
	}

	return &Storage{
		client: c,
		bucket: cfg.Bucket,
	}, nil
}

// Upload image to the minio bucket
func (s *Storage) Upload(ctx context.Context, r io.Reader, size int64, contentType string) (string, error) {
	id := uuid.New().String()
	if _, err := s.client.PutObject(
		ctx, s.bucket, id, r, size, minio.PutObjectOptions{
			ContentType: contentType,
		},
	); err != nil {
		return "", fmt.Errorf("unable to upload image: %w", err)
	}

	return id, nil
}
