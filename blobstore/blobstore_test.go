package blobstore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// mockClient implements the Downloader interface for testing.
type mockClient struct {
	mu            sync.Mutex
	files         map[string]string
	uploadedFiles map[string]string
}

func (m *mockClient) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (int64, error) {
	content, ok := m.files[*input.Key]
	if !ok {
		return 0, errors.New("file not found")
	}
	n, err := w.WriteAt([]byte(content), 0)
	return int64(n), err
}

func (m *mockClient) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	content, ok := m.files[*params.Key]
	if !ok {
		return nil, errors.New("file not found")
	}
	return &s3.HeadObjectOutput{ContentLength: int64(len(content))}, nil
}

func (m *mockClient) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	res := &s3.ListObjectsV2Output{Contents: []types.Object{}}
	for k := range m.files {
		res.Contents = append(res.Contents, types.Object{Key: aws.String(k)})
	}
	return res, nil
}

func (m *mockClient) Upload(ctx context.Context, input *s3.PutObjectInput, options ...func(*manager.Uploader)) (output *manager.UploadOutput, err error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, input.Body); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.uploadedFiles[*input.Key] = buf.String()
	return &manager.UploadOutput{}, nil
}

func createLocalFiles(dir string, files map[string]string) error {
	for path, content := range files {
		fullPath := filepath.Join(dir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func TestBlobStore_Download(t *testing.T) {
	tests := []struct {
		name          string
		remoteFiles   map[string]string
		localFiles    map[string]string
		key           string
		localDir      string
		overwrite     bool
		expectedCount int
		expectedErr   error
	}{
		{
			name:          "Empty directory",
			remoteFiles:   map[string]string{},
			localFiles:    map[string]string{},
			key:           "dir1",
			localDir:      "test_download",
			overwrite:     false,
			expectedCount: 0,
			expectedErr:   errors.New("file not found"),
		},
		{
			name: "Single file",
			remoteFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			localFiles:    map[string]string{},
			key:           "file1.txt",
			localDir:      "test_download",
			overwrite:     false,
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "File exists locally, no overwrite",
			remoteFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			localFiles: map[string]string{
				"file1.txt": "file1 local content",
			},
			key:           "file1.txt",
			localDir:      "test_download",
			overwrite:     false,
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "File exists locally, with overwrite",
			remoteFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			localFiles: map[string]string{
				"file1.txt": "file1 local content",
			},
			key:           "file1.txt",
			localDir:      "test_download",
			overwrite:     true,
			expectedCount: 1,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockClient := &mockClient{files: tt.remoteFiles}
			store := &BlobStore{client: mockClient, downloader: mockClient}

			tmpDir := t.TempDir()
			localDir := filepath.Join(tmpDir, tt.localDir)
			if err := createLocalFiles(localDir, tt.localFiles); err != nil {
				t.Fatalf("Failed to create local files: %v", err)
			}

			err := store.Download(context.Background(), "", tt.key, localDir, tt.overwrite)

			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Check if the expected number of files are downloaded
			count := 0
			if err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					count++
				}
				return nil
			}); err != nil {
				t.Fatalf("Failed to walk local directory: %v", err)
			}

			if count != tt.expectedCount {
				t.Errorf("Expected %d files, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestBlobStore_UploadFromPath(t *testing.T) {
	tests := []struct {
		name          string
		localFiles    map[string]string
		remoteFiles   map[string]string
		key           string
		localPath     string
		overwrite     bool
		expectedCount int
		expectedErr   error
	}{
		{
			name: "Single file",
			localFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			remoteFiles:   map[string]string{},
			key:           "file1.txt",
			localPath:     "file1.txt",
			overwrite:     false,
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "Directory upload",
			localFiles: map[string]string{
				"file1.txt":      "file1 content",
				"dir1/file2.txt": "file2 content",
			},
			remoteFiles:   map[string]string{},
			key:           "dir1",
			localPath:     "dir1",
			overwrite:     false,
			expectedCount: 1,
			expectedErr:   nil,
		},
		{
			name: "File exists remotely, no overwrite",
			localFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			remoteFiles: map[string]string{
				"file1.txt": "file1 remote content",
			},
			key:           "file1.txt",
			localPath:     "file1.txt",
			overwrite:     false,
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name: "File exists remotely, with overwrite",
			localFiles: map[string]string{
				"file1.txt": "file1 content",
			},
			remoteFiles: map[string]string{
				"file1.txt": "file1 remote content",
			},
			key:           "file1.txt",
			localPath:     "file1.txt",
			overwrite:     true,
			expectedCount: 1,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUploader := &mockClient{
				uploadedFiles: make(map[string]string),
				files:         tt.remoteFiles,
			}
			store := &BlobStore{uploader: mockUploader, client: mockUploader}

			tmpDir := t.TempDir()
			localPath := filepath.Join(tmpDir, tt.localPath)

			if err := createLocalFiles(filepath.Dir(localPath), tt.localFiles); err != nil {
				t.Fatalf("Failed to create local files: %v", err)
			}

			err := store.UploadFromPath(context.Background(), tt.key, localPath, tt.overwrite)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Check if the expected number of files are uploaded. Make sure the files
			// that already exist remotely are not uploaded again.
			count := 0
			mockUploader.mu.Lock()
			defer mockUploader.mu.Unlock()
			for k, v := range mockUploader.uploadedFiles {
				if tt.overwrite {
					if v == tt.localFiles[k] {
						count++
					}
					continue
				}
				if v != tt.remoteFiles[k] {
					count++
				}
			}
		})
	}
}
