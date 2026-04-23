package parquet

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/parquet-go/parquet-go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet/scheme"
	"github.com/reportportal/service-ingest/internal/model"
)

func (w *Writer) writeLogEvents(ow io.Writer, events []buffer.EventEnvelope) error {
	rows := make([]scheme.LogEvent, 0, len(events))

	for _, event := range events {
		var log model.Log
		if err := json.Unmarshal(event.Data, &log); err != nil {
			return fmt.Errorf("failed to unmarshal log: %w", err)
		}

		row := scheme.NewLogEvent(event, log)
		rows = append(rows, row)
	}

	pw := parquet.NewGenericWriter[scheme.LogEvent](ow, parquet.Compression(w.codec))
	if _, err := pw.Write(rows); err != nil {
		_ = pw.Close()
		return fmt.Errorf("failed to write log events: %w", err)
	}

	return pw.Close()
}
