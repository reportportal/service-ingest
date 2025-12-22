package parquet

import (
	"encoding/json"
	"fmt"

	"github.com/parquet-go/parquet-go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet/scheme"
	"github.com/reportportal/service-ingest/internal/model"
)

func (w *Writer) writeLaunchEvents(filename string, events []buffer.EventEnvelope) error {
	rows := make([]scheme.LaunchEvent, 0, len(events))

	for _, event := range events {
		var launch model.Launch
		if err := json.Unmarshal(event.Data, &launch); err != nil {
			return fmt.Errorf("failed to unmarshal launch: %w", err)
		}

		row := scheme.NewLaunchEvent(event, launch)
		rows = append(rows, row)
	}

	return parquet.WriteFile(
		filename,
		rows,
		parquet.Compression(w.getCompressionCodec()),
	)
}
