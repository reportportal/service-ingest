package processor

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type FileProcessor struct {
	buffer        buffer.FileBuffer
	catalogDir    string
	flushInterval time.Duration
	logger        *slog.Logger

	done chan struct{}
}

func NewFileProcessor(fileBuffer buffer.FileBuffer, dir string, flushInterval time.Duration) *FileProcessor {
	return &FileProcessor{
		buffer:        fileBuffer,
		catalogDir:    dir,
		flushInterval: flushInterval,
		logger:        slog.Default(),
		done:          make(chan struct{}),
	}
}

func (fp *FileProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(fp.flushInterval)
	defer ticker.Stop()
	defer close(fp.done)

	fp.logger.Info("starting file processor", "flush_interval", fp.flushInterval)

	for {
		select {
		case <-ctx.Done():
			fp.logger.Info("stopping file processor")
			return
		case <-ticker.C:
			if err := fp.Flush(ctx); err != nil {
				fp.logger.Warn("failed to flush files", "error", err)
			}
		}

	}
}

func (fp *FileProcessor) Done() <-chan struct{} {
	return fp.done
}

func (fp *FileProcessor) Flush(ctx context.Context) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	files, err := fp.buffer.List()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fp.logger.Debug("file buffer is empty, skipping batch")
		return nil
	}

	fp.logger.Debug("flushing file processor", "count", len(files))

	for _, f := range files {
		if err := fp.move(f); err != nil {
			fp.logger.Error("failed to move file", "file", f, "error", err)
			continue
		}
	}

	return nil
}

func (fp *FileProcessor) move(f string) error {
	dir, file := filepath.Split(f)
	dir = filepath.Clean(dir)
	destDir := filepath.Join(fp.catalogDir, dir)
	destFile := filepath.Join(destDir, file)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", destDir, err)
	}

	src, err := fp.buffer.Read(f)
	if err != nil {
		return fmt.Errorf("read file %s from buffer: %w", f, err)
	}
	defer src.Close()

	tmp, err := os.CreateTemp(destDir, "move-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, src); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("write to temp file: %w", err)
	}

	if err := os.Rename(tmp.Name(), destFile); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("rename file %s to %s: %w", tmp.Name(), destFile, err)
	}

	if err := fp.buffer.Delete(f); err != nil {
		return fmt.Errorf("delete buffer file %s: %w", f, err)
	}

	return nil
}
