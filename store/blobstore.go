package store

import (
	"context"
	"encoding/hex"
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

type Store interface {
	List(ctx context.Context, prefix string) ([]string, error)
	DownloadToPath(ctx context.Context, key, localDir string, localFileHashes map[string]string) error
	RemoteObjectMD5(ctx context.Context, key string) (string, error)
	Upload(ctx context.Context, key string, content io.Reader) error
	UploadFromPath(ctx context.Context, key, localPath string) error
}

type BlobStore struct {
	client     *s3.Client
	bucket     string
	downloader *manager.Downloader
	uploader   *manager.Uploader
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

func (b *BlobStore) DownloadToPath(ctx context.Context, key, localDir string, localFileHashes map[string]string) error {
	localPath := filepath.Join(localDir, filepath.FromSlash(key))
	dir := filepath.Dir(localPath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	remoteMd5, err := b.RemoteObjectMD5(ctx, key)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Error getting remote object MD5 for '%s/%s'", b.bucket, key)
	}

	// Check if the file exists locally with the same content
	existingLocalPath, exists := localFileHashes[remoteMd5]
	if exists {
		err := copyFile(existingLocalPath, localPath)
		if err != nil {
			return err
		}
		log.Ctx(ctx).Info().Msgf("Copied existing local file '%s' to '%s'", existingLocalPath, localPath)
		return nil
	}

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.GetObjectInput{
		Bucket: &b.bucket,
		Key:    &key,
	}

	_, err = b.downloader.Download(ctx, file, input)
	if err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msgf("Successfully downloaded '%s/%s' to '%s'", b.bucket, key, localPath)

	return nil
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

func (b *BlobStore) Upload(ctx context.Context, key string, content io.Reader) error {
	log.Ctx(ctx).Info().Msgf("Uploading '%s' to '%s/%s'", key, b.bucket, key)

	input := &s3.PutObjectInput{
		Bucket: &b.bucket,
		Key:    aws.String(key),
		Body:   content,
	}

	_, err := b.uploader.Upload(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlobStore) UploadFromPath(ctx context.Context, key, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if strings.HasPrefix(key, string(os.PathSeparator)) {
		key = strings.TrimPrefix(key, string(os.PathSeparator))
	}
	key = filepath.ToSlash(key)

	if err := b.Upload(ctx, key, file); err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msgf("Successfully uploaded '%s' to '%s/%s'", localPath, b.bucket, key)

	return nil
}

func NewBlobStore(bucket string, s3Cfg aws.Config) *BlobStore {
	client := s3.NewFromConfig(s3Cfg)
	return &BlobStore{
		client:     client,
		bucket:     bucket,
		downloader: manager.NewDownloader(client),
		uploader:   manager.NewUploader(client),
	}
}
