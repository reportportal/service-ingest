package parquet

import (
	"fmt"
	"path"

	opendal "github.com/apache/opendal/bindings/go"
	"github.com/parquet-go/parquet-go/compress"
	"github.com/parquet-go/parquet-go/compress/gzip"
	"github.com/parquet-go/parquet-go/compress/snappy"
	"github.com/parquet-go/parquet-go/compress/uncompressed"
	"github.com/parquet-go/parquet-go/compress/zstd"
	"github.com/reportportal/service-ingest/internal/data/buffer"
)

type Writer struct {
	codec    compress.Codec
	operator *opendal.Operator
}

func NewWriter(compression string, operator *opendal.Operator) *Writer {
	return &Writer{
		codec:    getCompressionCodec(compression),
		operator: operator,
	}
}

func (w *Writer) Write(entityType buffer.EntityType, dir string, events []buffer.EventEnvelope) error {
	if len(events) == 0 {
		return nil
	}

	// TODO: Implement partitioning logic if needed
	parquetPath := path.Join(dir, "part-00000.parquet")

	if err := w.writeParquet(parquetPath, entityType, events); err != nil {
		return fmt.Errorf("failed to write parquet file: %w", err)
	}

	if err := w.operator.Write(path.Join(dir, "_SUCCESS"), []byte{}); err != nil {
		return fmt.Errorf("failed to create _SUCCESS marker: %w", err)
	}

	return nil
}

func (w *Writer) writeParquet(parquetPath string, entityType buffer.EntityType, events []buffer.EventEnvelope) (err error) {
	writer, err := w.operator.Writer(parquetPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := writer.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close parquet writer: %w", closeErr)
		}
	}()

	switch entityType {
	case buffer.EntityTypeLaunch:
		return w.writeLaunchEvents(writer, events)
	case buffer.EntityTypeItem:
		return w.writeItemEvents(writer, events)
	case buffer.EntityTypeLog:
		return w.writeLogEvents(writer, events)
	default:
		return fmt.Errorf("unknown entity type: %s", entityType)
	}
}

func getCompressionCodec(codec string) compress.Codec {
	switch codec {
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
