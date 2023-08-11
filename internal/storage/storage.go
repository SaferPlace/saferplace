// Copyright 2023 SaferPlace

package storage

import (
	"context"
	"io"
)

// Storage allows to upload the image
type Storage interface {
	// Upload takes in the reader from which it reads from to get the image and returns the
	// reference which can uniquely identify the image, or an error if there was a problem uploading
	// to the bucket.
	Upload(ctx context.Context, r io.Reader, size int64, contentType string) (string, error)
}
