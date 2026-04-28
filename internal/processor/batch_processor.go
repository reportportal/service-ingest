package processor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/catalog"
	"github.com/reportportal/service-ingest/internal/data/parquet"
)

type BatchProcessor struct {
	buffer        buffer.Buffer
	writer        *parquet.Writer
	flushInterval time.Duration
	readLimit     int
	logger        *slog.Logger

	done chan struct{}
}

type BatchProcessorOptions struct {
	Buffer        buffer.Buffer
	Writer        *parquet.Writer
	FlushInterval time.Duration
	ReadLimit     int
	Logger        *slog.Logger
}

func NewBatchProcessor(opts BatchProcessorOptions) *BatchProcessor {
	return &BatchProcessor{
		buffer:        opts.Buffer,
		writer:        opts.Writer,
		flushInterval: opts.FlushInterval,
		readLimit:     opts.ReadLimit,
		logger:        opts.Logger,
		done:          make(chan struct{}),
	}
}

func (bp *BatchProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()
	defer close(bp.done)

	bp.logger.Info("batch processor started",
		"flush_interval", bp.flushInterval,
		"read_limit", bp.readLimit,
	)

	for {
		select {
		case <-ctx.Done():
			bp.logger.Info("batch processor stopped")
			return
		case <-ticker.C:
			if err := bp.processBatch(ctx); err != nil {
				bp.logger.Error("failed to process batch", "error", err)
			}
		}
	}
}

func (bp *BatchProcessor) Done() <-chan struct{} {
	return bp.done
}

func (bp *BatchProcessor) processBatch(ctx context.Context) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	counter, err := bp.buffer.Size(ctx)
	if err != nil {
		return fmt.Errorf("failed to get buffer size: %w", err)
	}

	if counter == 0 {
		bp.logger.Debug("buffer is empty, skipping batch")
		return nil
	}

	bp.logger.Debug("processing batch", "count", counter)

	events, err := bp.buffer.Read(ctx, bp.readLimit)
	if err != nil {
		return fmt.Errorf("failed to read from buffer: %w", err)
	}

	if len(events) == 0 {
		bp.logger.Debug("no events to process")
		return nil
	}

	bp.logger.Debug("read events from buffer", "count", len(events))

	batchID := fmt.Sprintf("%d", time.Now().UnixMilli())
	groups := bp.group(events)

	bp.logger.Debug("grouped events", "groups", len(groups))

	for key, group := range groups {
		if err := bp.write(key, batchID, group); err != nil {
			bp.logger.Error("failed to write group", "groups", key, "error", err, "batch_id", batchID)
			return fmt.Errorf("failed to write partition %v: %w", key, err)
		}

		if err := bp.buffer.Ack(ctx, group); err != nil {
			return fmt.Errorf("failed to ack events: %w", err)
		}
	}

	bp.logger.Debug("batch processed successfully", "events", len(events), "groups", len(groups))

	return nil
}

func (bp *BatchProcessor) group(events []buffer.EventEnvelope) map[groupKey][]buffer.EventEnvelope {
	groups := make(map[groupKey][]buffer.EventEnvelope)

	for _, event := range events {
		key := newGroupKey(event)
		groups[key] = append(groups[key], event)
	}

	return groups
}

func (bp *BatchProcessor) write(key groupKey, batchID string, events []buffer.EventEnvelope) error {
	path := catalog.BuildPath(
		key.ProjectKey,
		key.LaunchUUID,
		string(key.EntityType),
		key.EventDate,
		batchID,
	)

	bp.logger.Debug("writing group", "path", path, "events", len(events))

	if err := bp.writer.Write(key.EntityType, path, events); err != nil {
		return fmt.Errorf("failed to write partition: %w", err)
	}

	return nil
}
