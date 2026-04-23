package processor

import (
	"bytes"
	"crypto/sha256"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/opendal-go-services/fs"
	opendal "github.com/apache/opendal/bindings/go"
	"github.com/stretchr/testify/require"
)

func TestFileProcessor_upload_LargeFile(t *testing.T) {
	const fileSize = 32 << 20 // 32 MiB

	tmpDir := t.TempDir()
	op, err := opendal.NewOperator(fs.Scheme, opendal.OperatorOptions{"root": tmpDir})
	require.NoError(t, err)
	t.Cleanup(op.Close)

	fp := &FileProcessor{
		operator: op,
		logger:   slog.Default(),
	}

	src := make([]byte, fileSize)
	_, err = rand.New(rand.NewSource(1)).Read(src)
	require.NoError(t, err)

	const destPath = "attachments/large.bin"
	require.NoError(t, fp.upload(destPath, bytes.NewReader(src)))

	f, err := os.Open(filepath.Join(tmpDir, destPath))
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	info, err := f.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(fileSize), info.Size(), "uploaded size should match source")

	srcHash := sha256.Sum256(src)
	dstHasher := sha256.New()
	_, err = io.Copy(dstHasher, f)
	require.NoError(t, err)
	require.Equal(t, srcHash[:], dstHasher.Sum(nil), "uploaded content should match source byte-for-byte")
}
