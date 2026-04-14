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
	Dir string
}

func NewFileBuffer(dir string) FileBuffer {
	return FileBuffer{dir}
}

func (fb *FileBuffer) Save(path string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open multipart file: %w", err)
	}
	defer src.Close()

	hasher := sha256.New()

	fullPath := filepath.Join(fb.Dir, path)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", fmt.Errorf("create buffer directory: %w", err)
	}

	tmpDir := filepath.Join(fb.Dir, "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("create temp buffer directory: %w", err)
	}

	tmp, err := os.CreateTemp(tmpDir, "upload-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, io.TeeReader(src, hasher)); err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("write file: %w", err)
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	if err := os.Rename(tmp.Name(), filepath.Join(fullPath, hash)); err != nil {
		_ = os.Remove(tmp.Name())
		return "", fmt.Errorf("rename to %s: %w", hash, err)
	}

	return hash, nil
}

// List returns list of relative file paths.
func (fb *FileBuffer) List() (files []string, err error) {
	err = filepath.WalkDir(fb.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == "tmp" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(fb.Dir, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		files = append(files, rel)
		return nil
	})

	return files, err
}

// Read takes a relative path to file to read it.
func (fb *FileBuffer) Read(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(fb.Dir, path))
}

// Delete take relative path to file to delete it and catalog if it's empty.
func (fb *FileBuffer) Delete(path string) error {
	if err := os.Remove(filepath.Join(fb.Dir, path)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file %s: %w", path, err)
	}

	dir := filepath.Dir(filepath.Join(fb.Dir, path))
	for dir != fb.Dir {
		if err := os.Remove(dir); err != nil {
			break
		}
		dir = filepath.Dir(dir)
	}

	return nil
}
