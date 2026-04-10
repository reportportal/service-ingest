package buffer

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMultipartFileHeader(t *testing.T, content []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.txt"`)
	h.Set("Content-Type", "application/octet-stream")

	part, err := w.CreatePart(h)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	w.Close()

	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(int64(len(content) + 1024))
	require.NoError(t, err)

	files := form.File["file"]
	require.NotEmpty(t, files)
	return files[0]
}

func expectedHash(content []byte) string {
	h := sha256.Sum256(content)
	return fmt.Sprintf("%x", h[:])
}

func TestFileBuffer_Save(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
	}{
		{
			name:    "saves file and returns sha256 hash",
			content: []byte("hello world"),
		},
		{
			name:    "saves empty file",
			content: []byte(""),
		},
		{
			name:    "saves binary content",
			content: []byte{0x00, 0xFF, 0xAB, 0xCD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			fb := NewFileBuffer(dir)
			path := filepath.Join("test-project", uuid.New().String())

			fh := createMultipartFileHeader(t, tt.content)
			hash, err := fb.Save(path, fh)
			require.NoError(t, err)

			assert.Equal(t, expectedHash(tt.content), hash)

			got, err := os.ReadFile(filepath.Join(dir, path, hash))
			require.NoError(t, err)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestFileBuffer_Save_DuplicateContent(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)
	path := filepath.Join("test-project", uuid.New().String())
	content := []byte("duplicate")

	fh1 := createMultipartFileHeader(t, content)
	hash1, err := fb.Save(path, fh1)
	require.NoError(t, err)

	fh2 := createMultipartFileHeader(t, content)
	hash2, err := fb.Save(path, fh2)
	require.NoError(t, err)

	assert.Equal(t, hash1, hash2)
}

func TestFileBuffer_Read(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)
	path := filepath.Join("test-project", uuid.New().String())
	content := []byte("read me")

	fh := createMultipartFileHeader(t, content)
	hash, err := fb.Save(path, fh)
	require.NoError(t, err)

	rc, err := fb.Read(path, hash)
	require.NoError(t, err)
	defer rc.Close()

	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, content, got)
}

func TestFileBuffer_Read_NotFound(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)

	_, err := fb.Read("", "nonexistent")
	assert.Error(t, err)
}

func TestFileBuffer_Delete(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)
	path := filepath.Join("test-project", uuid.New().String())
	content := []byte("delete me")

	fh := createMultipartFileHeader(t, content)
	hash, err := fb.Save(path, fh)
	require.NoError(t, err)

	require.NoError(t, fb.Delete(path, hash))

	_, err = os.Stat(filepath.Join(dir, hash))
	assert.True(t, os.IsNotExist(err))
}

func TestFileBuffer_Delete_NotFound(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)

	assert.NoError(t, fb.Delete("", "nonexistent"))
}

func TestFileBuffer_List(t *testing.T) {
	t.Run("empty buffer", func(t *testing.T) {
		dir := t.TempDir()
		fb := NewFileBuffer(dir)

		files, err := fb.List()
		require.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("lists files across projects and launches", func(t *testing.T) {
		dir := t.TempDir()
		fb := NewFileBuffer(dir)

		fh1 := createMultipartFileHeader(t, []byte("file1"))
		hash1, err := fb.Save(filepath.Join("project-a", "launch-1"), fh1)
		require.NoError(t, err)

		fh2 := createMultipartFileHeader(t, []byte("file2"))
		hash2, err := fb.Save(filepath.Join("project-a", "launch-2"), fh2)
		require.NoError(t, err)

		fh3 := createMultipartFileHeader(t, []byte("file3"))
		hash3, err := fb.Save(filepath.Join("project-b", "launch-3"), fh3)
		require.NoError(t, err)

		files, err := fb.List()
		require.NoError(t, err)
		assert.Len(t, files, 3)

		expected := []string{
			filepath.Join("project-a", "launch-1", hash1),
			filepath.Join("project-a", "launch-2", hash2),
			filepath.Join("project-b", "launch-3", hash3),
		}
		assert.ElementsMatch(t, expected, files)
	})
}

func TestFileBuffer_SaveReadDelete_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	fb := NewFileBuffer(dir)
	path := filepath.Join("test-project", uuid.New().String())
	content := []byte("full roundtrip")

	fh := createMultipartFileHeader(t, content)
	hash, err := fb.Save(path, fh)
	require.NoError(t, err)

	rc, err := fb.Read(path, hash)
	require.NoError(t, err)
	got, err := io.ReadAll(rc)
	require.NoError(t, err)
	rc.Close()

	assert.Equal(t, content, got)

	require.NoError(t, fb.Delete(path, hash))

	_, err = fb.Read(path, hash)
	assert.Error(t, err)
}
