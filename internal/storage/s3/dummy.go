package s3

import (
	"io"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"golang.org/x/net/context"
)

type S3DummyRepo struct {
}

func NewS3DummyRepo() *S3DummyRepo {
	return &S3DummyRepo{}
}

func (s *S3DummyRepo) BucketExists(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *S3DummyRepo) PutObject(ctx context.Context, objName, contentType string, contentLength int64, r io.Reader) error {
	return nil
}

func (s *S3DummyRepo) GetObject(ctx context.Context, objName string) (io.ReadSeekCloser, error) {
	return nil, storage.ErrS3ServiceDisabled
}
