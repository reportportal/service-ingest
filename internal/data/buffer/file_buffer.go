package buffer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileBuffer struct {
	Path string
}

func NewFileBuffer(path string) FileBuffer {
	return FileBuffer{Path: path}
}

func (fb *FileBuffer) Save(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open multipart file: %w", err)
	}
	defer src.Close()

	hasher := sha256.New()

	tmp, err := os.CreateTemp(fb.Path, "upload-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, io.TeeReader(src, hasher)); err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("write file: %w", err)
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	dest := filepath.Join(fb.Path, hash)

	if err := os.Rename(tmp.Name(), dest); err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("rename to %s: %w", hash, err)
	}

	return hash, nil
}

func (fb *FileBuffer) Read(hash string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fb.Path, hash))
}

func (fb *FileBuffer) Delete(hash string) error {
	if err := os.Remove(filepath.Join(fb.Path, hash)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file %s: %w", hash, err)
	}
	return nil
}
