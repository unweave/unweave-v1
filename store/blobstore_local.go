package store

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type LocalBlobStore struct {
	rootDir string
}

func (l *LocalBlobStore) List(ctx context.Context, prefix string) ([]string, error) {
	var objectKeys []string

	dir := filepath.Join(l.rootDir, prefix)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			objectKeys = append(objectKeys, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return objectKeys, nil
}

func (l *LocalBlobStore) Download(ctx context.Context, key, localDir string, overwrite bool) error {
	localPath := filepath.Join(localDir, filepath.FromSlash(key))
	remotePath := filepath.Join(l.rootDir, key)
	dir := filepath.Dir(localPath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	err = copyFile(remotePath, localPath)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalBlobStore) RemoteObjectMD5(ctx context.Context, key string) (string, error) {
	src := filepath.Join(l.rootDir, key)

	file, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	remoteMd5 := hex.EncodeToString(hashInBytes)

	return remoteMd5, nil
}

func (l *LocalBlobStore) Upload(ctx context.Context, key string, content io.Reader, overwrite bool) error {
	dst := filepath.Join(l.rootDir, key)

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, content)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalBlobStore) UploadFromPath(ctx context.Context, key, localPath string, overwrite bool) error {
	relPath, err := filepath.Rel(localPath, l.rootDir)
	if err != nil {
		return err
	}

	dst := filepath.Join(key, relPath)

	err = copyFile(localPath, dst)
	if err != nil {
		return err
	}

	return nil
}

func NewLocalBlobStore(rootDir string) *LocalBlobStore {
	return &LocalBlobStore{
		rootDir: rootDir,
	}
}
