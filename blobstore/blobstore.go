package blobstore

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

type Downloader interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)
}

type Uploader interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, options ...func(*manager.Uploader)) (output *manager.UploadOutput, err error)
}

// S3Client is an interface defining a subset of the S3 client methods used by BlobStore.
type S3Client interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type Store interface {
	Download(ctx context.Context, remoteDir, remoteKey, localDir string, overwrite bool) error
	List(ctx context.Context, prefix string) ([]string, error)
	RemoteObjectMD5(ctx context.Context, key string) (string, error)
	Upload(ctx context.Context, key string, content io.Reader, overwrite bool) error
	UploadFromPath(ctx context.Context, key, localPath string, overwrite bool) error
}

type BlobStore struct {
	client     S3Client
	bucket     string
	downloader Downloader
	uploader   Uploader
}

func (b *BlobStore) List(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.bucket),
		Prefix: aws.String(prefix),
	}

	var objectKeys []string

	for {
		output, err := b.client.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, object := range output.Contents {
			objectKeys = append(objectKeys, *object.Key)
		}

		if !output.IsTruncated {
			break
		}
		input.ContinuationToken = output.NextContinuationToken
	}
	return objectKeys, nil
}

func (b *BlobStore) Download(ctx context.Context, remoteDir, remoteKey, localDir string, overwrite bool) error {
	isDir, err := b.isRemoteKeyDir(ctx, remoteKey)
	if err != nil {
		return err
	}
	if isDir {
		return b.downloadDirectory(ctx, remoteKey, localDir, overwrite)
	}

	log.Info().Msgf("Downloading '%s/%s' to '%s'", b.bucket, remoteKey, localDir)

	rel, err := filepath.Rel(remoteDir, remoteKey)
	if err != nil {
		return fmt.Errorf("invalid remoteDir %q for remoteKey %q: %v", remoteDir, remoteKey, err)
	}

	path := filepath.Join(localDir, rel)
	dir := filepath.Dir(path)

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if !overwrite {
		// Check if file with the name already exists
		if _, err := os.Stat(path); err == nil {
			log.Info().Msgf("File '%s' already exists, skipping download", path)
			return nil
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.GetObjectInput{
		Bucket: &b.bucket,
		Key:    &remoteKey,
	}

	_, err = b.downloader.Download(ctx, file, input)
	if err != nil {
		if e := os.Remove(path); e != nil {
			log.Error().Err(e).Msgf("Failed to remove file '%s'", path)
		}
		return err
	}
	log.Info().Msgf("Successfully downloaded '%s/%s' to '%s'", b.bucket, remoteKey, path)

	return nil
}

func (b *BlobStore) downloadDirectory(ctx context.Context, remoteDir, localDir string, overwrite bool) error {
	log.Info().Msgf("Downloading directory '%s/%s' to '%s'", b.bucket, remoteDir, localDir)

	if err := os.MkdirAll(localDir, os.ModePerm); err != nil {
		return err
	}

	keys, err := b.List(ctx, remoteDir)
	if err != nil {
		return fmt.Errorf("failed to list objects, %v", err)
	}

	for _, key := range keys {
		log.Info().Msgf("Downloading '%s/%s' to '%s'", b.bucket, key, localDir)
		if err := b.Download(ctx, remoteDir, key, localDir, overwrite); err != nil {
			return err
		}
	}

	return nil
}

func (b *BlobStore) isRemoteKeyDir(ctx context.Context, key string) (bool, error) {
	if strings.HasSuffix(key, "/") {
		return true, nil
	}

	// It might be that the key is a directory, but the user forgot to add the trailing
	// slash. Let's make sure
	keys, err := b.List(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to list objects at key %q, %v", key, err)
	}
	if len(keys) > 1 {
		return true, nil
	}
	return false, nil
}

func (b *BlobStore) RemoteObjectMD5(ctx context.Context, key string) (string, error) {
	headInput := &s3.HeadObjectInput{
		Bucket: &b.bucket,
		Key:    &key,
	}
	headOutput, err := b.client.HeadObject(ctx, headInput)
	if err != nil {
		return "", err
	}

	remoteMd5 := ""
	if headOutput.ETag != nil {
		etag := *headOutput.ETag
		remoteMd5 = etag[1 : len(etag)-1]
	}
	remoteMd5 = hex.EncodeToString([]byte(remoteMd5))
	return remoteMd5, nil
}

func (b *BlobStore) Upload(ctx context.Context, key string, content io.Reader, overwrite bool) error {
	log.Info().Msgf("Uploading '%s' to '%s/%s'", key, b.bucket, key)

	input := &s3.PutObjectInput{
		Bucket: &b.bucket,
		Key:    aws.String(key),
		Body:   content,
	}

	if !overwrite {
		existing, err := b.List(ctx, key)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			log.Info().Msgf("File '%s' already exists, skipping upload", key)
			return nil
		}
	}

	_, err := b.uploader.Upload(ctx, input)
	if err != nil {
		return err
	}
	log.Info().Msgf("Successfully uploaded file to '%s/%s'", b.bucket, key)

	return nil
}

func (b *BlobStore) uploadDirectory(ctx context.Context, key, localDir string, overwrite bool) error {
	log.Info().Msgf("Uploading directory '%s' to '%s/%s'", localDir, b.bucket, key)

	keys, err := b.List(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to list objects, %v", err)
	}

	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !overwrite {
			// Check if file with the name already exists
			for _, k := range keys {
				if k == relPath {
					log.Info().Msgf("File '%s' already exists, skipping upload", relPath)
					return nil
				}
			}
		}

		go func() {
			file, err := os.Open(path)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to open file '%s'", relPath)
				return
			}
			defer file.Close()
			if e := b.Upload(ctx, relPath, file, false); e != nil {
				log.Error().Err(e).Msgf("Failed to upload file '%s'", relPath)
			}
		}()

		return nil
	})

	return err
}

func (b *BlobStore) UploadFromPath(ctx context.Context, key, localPath string, overwrite bool) error {
	stat, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path '%s' does not exist", localPath)
	}
	if stat.IsDir() {
		return b.uploadDirectory(ctx, key, localPath, overwrite)
	}
	log.Info().Msgf("Uploading '%s' to '%s/%s'", localPath, b.bucket, key)

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if strings.HasPrefix(key, string(os.PathSeparator)) {
		key = strings.TrimPrefix(key, string(os.PathSeparator))
	}
	key = filepath.ToSlash(key)

	if err := b.Upload(ctx, key, file, overwrite); err != nil {
		return err
	}

	return nil
}

func NewBlobStore(bucket string, s3Cfg aws.Config) *BlobStore {
	client := s3.NewFromConfig(s3Cfg)
	return &BlobStore{
		client: client,
		bucket: bucket,
		downloader: manager.NewDownloader(client, func(d *manager.Downloader) {
			d.PartSize = 64 * 1024 * 1024
			d.Concurrency = 10
		}),
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 64 * 1024 * 1024
			u.Concurrency = 10
		}),
	}
}
