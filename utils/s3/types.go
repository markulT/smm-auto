package s3

import (
	"context"
	"io"
)

type StorageService interface {
	Load(context.Context,string, io.Reader ) error
	Delete(string) error
	Get(string) (io.Reader, error)
}
