package parquet

import (
	"encoding/json"
	"fmt"

	"github.com/parquet-go/parquet-go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet/scheme"
	"github.com/reportportal/service-ingest/internal/model"
)

func (w *Writer) writeItemEvents(filename string, events []buffer.EventEnvelope) error {
	rows := make([]scheme.ItemEvent, 0, len(events))

	for _, event := range events {
		var item model.Item
		if err := json.Unmarshal(event.Data, &item); err != nil {
			return fmt.Errorf("failed to unmarshal item: %w", err)
		}

		row := scheme.NewItemEvent(event, item)
		rows = append(rows, row)
	}

	return parquet.WriteFile(
		filename,
		rows,
		parquet.Compression(w.getCompressionCodec()),
	)
}
