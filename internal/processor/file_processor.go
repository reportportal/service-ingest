package processor

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	opendal "github.com/apache/opendal/bindings/go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type FileProcessor struct {
	buffer        buffer.FileBuffer
	flushInterval time.Duration
	logger        *slog.Logger
	operator      *opendal.Operator

	done chan struct{}
}

func NewFileProcessor(fileBuffer buffer.FileBuffer, flushInterval time.Duration, operator *opendal.Operator) *FileProcessor {
	return &FileProcessor{
		buffer:        fileBuffer,
		flushInterval: flushInterval,
		logger:        slog.Default(),
		operator:      operator,
		done:          make(chan struct{}),
	}
}

func (fp *FileProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(fp.flushInterval)
	defer ticker.Stop()
	defer close(fp.done)

	fp.logger.Info("starting file processor", "file_buffer_dir", fp.buffer.Dir, "flush_interval", fp.flushInterval)

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
	src, err := fp.buffer.Read(f)
	if err != nil {
		return fmt.Errorf("read file %s from buffer: %w", f, err)
	}
	defer src.Close()

	if err := fp.upload(f, src); err != nil {
		return err
	}

	if err := fp.buffer.Delete(f); err != nil {
		return fmt.Errorf("delete buffer file %s: %w", f, err)
	}

	return nil
}

func (fp *FileProcessor) upload(f string, src io.Reader) (err error) {
	path := filepath.ToSlash(f)

	writer, err := fp.operator.Writer(path)
	if err != nil {
		return fmt.Errorf("create file %s writer: %w", f, err)
	}
	defer func() {
		if closeErr := writer.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file writer: %w", closeErr)
		}
	}()

	if _, err := io.Copy(writer, src); err != nil {
		return fmt.Errorf("copy file %s to %s: %w", f, path, err)
	}
	return nil
}
