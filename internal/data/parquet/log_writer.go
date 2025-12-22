package parquet

import (
	"encoding/json"
	"fmt"

	"github.com/parquet-go/parquet-go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet/scheme"
	"github.com/reportportal/service-ingest/internal/model"
)

func (w *Writer) writeLogEvents(filename string, events []buffer.EventEnvelope) error {
	rows := make([]scheme.LogEvent, 0, len(events))

	for _, event := range events {
		var log model.Log
		if err := json.Unmarshal(event.Data, &log); err != nil {
			return fmt.Errorf("failed to unmarshal log: %w", err)
		}

		row := scheme.NewLogEvent(event, log)
		rows = append(rows, row)
	}

	return parquet.WriteFile(
		filename,
		rows,
		parquet.Compression(w.getCompressionCodec()),
	)
}
