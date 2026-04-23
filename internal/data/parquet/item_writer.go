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

func (w *Writer) writeItemEvents(ow io.Writer, events []buffer.EventEnvelope) error {
	rows := make([]scheme.ItemEvent, 0, len(events))

	for _, event := range events {
		var item model.Item
		if err := json.Unmarshal(event.Data, &item); err != nil {
			return fmt.Errorf("failed to unmarshal item: %w", err)
		}

		row := scheme.NewItemEvent(event, item)
		rows = append(rows, row)
	}

	pw := parquet.NewGenericWriter[scheme.ItemEvent](ow, parquet.Compression(w.codec))
	if _, err := pw.Write(rows); err != nil {
		_ = pw.Close()
		return fmt.Errorf("failed to write item events: %w", err)
	}

	return pw.Close()
}
