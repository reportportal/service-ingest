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

func (w *Writer) writeLaunchEvents(ow io.Writer, events []buffer.EventEnvelope) error {
	rows := make([]scheme.LaunchEvent, 0, len(events))

	for _, event := range events {
		var launch model.Launch
		if err := json.Unmarshal(event.Data, &launch); err != nil {
			return fmt.Errorf("failed to unmarshal launch: %w", err)
		}

		row := scheme.NewLaunchEvent(event, launch)
		rows = append(rows, row)
	}

	pw := parquet.NewGenericWriter[scheme.LaunchEvent](ow, parquet.Compression(w.codec))
	if _, err := pw.Write(rows); err != nil {
		_ = pw.Close()
		return fmt.Errorf("failed to write launch events: %w", err)
	}

	return pw.Close()
}
