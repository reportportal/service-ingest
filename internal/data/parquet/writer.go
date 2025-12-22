package parquet

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/parquet-go/parquet-go/compress"
	"github.com/parquet-go/parquet-go/compress/gzip"
	"github.com/parquet-go/parquet-go/compress/snappy"
	"github.com/parquet-go/parquet-go/compress/uncompressed"
	"github.com/parquet-go/parquet-go/compress/zstd"
	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type Writer struct {
	basePath    string
	compression string
}

func NewWriter(basePath string, compression string) *Writer {
	return &Writer{
		basePath:    basePath,
		compression: compression,
	}
}

func (w *Writer) WritePartition(entityType buffer.EntityType, path string, events []buffer.EventEnvelope) error {
	if len(events) == 0 {
		return nil
	}

	fullPath := filepath.Join(w.basePath, path)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
	}

	// TODO: Implement partitioning logic if needed
	parquetFile := filepath.Join(fullPath, "part-00000.parquet")

	switch entityType {
	case buffer.EntityTypeLaunch:
		if err := w.writeLaunchEvents(parquetFile, events); err != nil {
			return err
		}
	case buffer.EntityTypeItem:
		if err := w.writeItemEvents(parquetFile, events); err != nil {
			return err
		}
	case buffer.EntityTypeLog:
		if err := w.writeLogEvents(parquetFile, events); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown entity type: %s", entityType)
	}

	successFile := filepath.Join(fullPath, "_SUCCESS")
	if err := os.WriteFile(successFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create _SUCCESS marker: %w", err)
	}

	return nil
}

func (w *Writer) getCompressionCodec() compress.Codec {
	switch w.compression {
	case "snappy":
		return &snappy.Codec{}
	case "gzip":
		return &gzip.Codec{}
	case "zstd":
		return &zstd.Codec{}
	case "uncompressed":
		return &uncompressed.Codec{}
	default:
		return &snappy.Codec{}
	}
}
