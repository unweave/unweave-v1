package tools

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src io.Reader, dest string) ([]string, error) {
	//nolint:prealloc
	var filenames []string

	readySrc, n, err := readerAt(src)
	if err != nil {
		return filenames, err
	}

	r, err := zip.NewReader(readySrc, n)
	if err != nil {
		return filenames, fmt.Errorf("zip.NewReader: %w", err)
	}

	for _, file := range r.File {
		// Store filename/path for returning and using later on
		fpath, err := securePath(dest, file.Name)
		if err != nil {
			return filenames, fmt.Errorf("securePath: %w", err)
		}

		filenames = append(filenames, fpath)

		outFile, err := prepareOutputFile(file, fpath)
		if errors.Is(err, errSkipDir) {
			continue
		} else if err != nil {
			return filenames, fmt.Errorf("prepareOutputFile: %w", err)
		}

		fileReader, err := file.Open()
		if err != nil {
			return filenames, fmt.Errorf("file.Open: %w", err)
		}

		// Protect against zip bombs. Max file size is 1GB

		_, err = io.Copy(outFile, io.LimitReader(fileReader, 1_000_000_000))

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		fileReader.Close()

		if err != nil {
			return filenames, fmt.Errorf("io.Copy: %w", err)
		}
	}

	return filenames, nil
}

type readerAtLen interface {
	io.ReaderAt
	Len() int
}

func readerAt(reader io.Reader) (io.ReaderAt, int64, error) {
	if r, ok := reader.(readerAtLen); ok {
		return r, int64(r.Len()), nil
	}

	var buf bytes.Buffer

	_, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, 0, fmt.Errorf("io.Copy from zip reader: %w", err)
	}

	r := bytes.NewReader(buf.Bytes())

	return r, int64(r.Len()), nil
}

func securePath(dest, name string) (string, error) {
	fpath := filepath.Join(dest, name)

	// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
	if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s: illegal file path", fpath)
	}

	return fpath, nil
}

var errSkipDir = fmt.Errorf("skip directory")

func prepareOutputFile(file *zip.File, fpath string) (io.WriteCloser, error) {
	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("os.MkdirAll: %w", err)
	}

	if file.FileInfo().IsDir() {
		return nil, errSkipDir
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	return outFile, nil
}
