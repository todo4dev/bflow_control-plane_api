package storage

import (
	"context"
	"io"
)

type WriteInput struct {
	FilePath string
	Stream   io.Reader
	MimeType string
}

type ReadResult struct {
	MimeType string
	Size     int64
	Stream   io.ReadCloser
}

type IStorageAdapter interface {
	Ping(ctx context.Context) error
	Write(ctx context.Context, input WriteInput) (string, error)
	Read(ctx context.Context, filePath string, digest string) (*ReadResult, error)
	Delete(ctx context.Context, filePath string, digest string) error
}
