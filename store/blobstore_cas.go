package store

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// CASStore is a store that compares files by their content before syncing them. If the
// content is the same, the file is not copied locally to the new filename.
type CASStore struct {
	client     *s3.Client
	bucket     string
	downloader *manager.Downloader
	uploader   *manager.Uploader
	blobstore  *BlobStore
}

func (s *CASStore) List(ctx context.Context, prefix string) ([]string, error) {
	return s.blobstore.List(ctx, prefix)
}

func (s *CASStore) Download(ctx context.Context, key, localDir string, overwrite bool) error {
	return nil
}

func (s *CASStore) RemoteObjectMD5(ctx context.Context, key string) (string, error) {
	return s.blobstore.RemoteObjectMD5(ctx, key)
}

func (s *CASStore) Upload(ctx context.Context, key string, content io.Reader, overwrite bool) error {
	return nil
}

func (s *CASStore) UploadFromPath(ctx context.Context, key, localPath string, overwrite bool) error {
	return nil
}
