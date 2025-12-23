package processor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/catalog"
	"github.com/reportportal/service-ingest/internal/data/parquet"
)

type BatchProcessor struct {
	buffer        buffer.Buffer
	writer        *parquet.Writer
	flushInterval time.Duration
	batchWindow   time.Duration
	readLimit     int
	logger        *slog.Logger

	//TODO: move to buffer for deterministic batch IDs across restarts
	mu     sync.Mutex // guards seqMap
	seqMap map[partitionKey]int

	done chan struct{}
}

type BatchProcessorOptions struct {
	Buffer        buffer.Buffer
	Writer        *parquet.Writer
	FlushInterval time.Duration
	BatchWindow   time.Duration
	ReadLimit     int
	Logger        *slog.Logger
}

func NewBatchProcessor(opts BatchProcessorOptions) *BatchProcessor {
	return &BatchProcessor{
		buffer:        opts.Buffer,
		writer:        opts.Writer,
		flushInterval: opts.FlushInterval,
		batchWindow:   opts.BatchWindow,
		readLimit:     opts.ReadLimit,
		logger:        opts.Logger,
		seqMap:        make(map[partitionKey]int),
		done:          make(chan struct{}),
	}
}

func (bp *BatchProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()
	defer close(bp.done)

	bp.logger.Info("batch processor started",
		"flush_interval", bp.flushInterval,
		"batch_window", bp.batchWindow,
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

func (bp *BatchProcessor) processBatch(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	counter, err := bp.buffer.Size(ctx)
	if err != nil {
		return fmt.Errorf("failed to get buffer size: %w", err)
	}

	if counter.Items == 0 {
		bp.logger.Debug("buffer is empty, skipping batch")
		return nil
	}

	bp.logger.Debug("processing batch", "count", counter.Items)

	events, err := bp.buffer.Read(ctx, bp.readLimit)
	if err != nil {
		return fmt.Errorf("failed to read from buffer: %w", err)
	}

	if len(events) == 0 {
		bp.logger.Debug("no events to process")
		return nil
	}

	bp.logger.Debug("read events from buffer", "count", len(events))

	partitions := bp.groupByPartition(events)

	bp.logger.Info("grouped events into partitions", "partitions", len(partitions))

	for key, partition := range partitions {
		batchID := bp.generateBatchID(key)

		if err := bp.writePartition(key, batchID, partition); err != nil {
			bp.logger.Error("failed to write partition", "partition", key, "error", err, "batch_id", batchID)

			if releaseErr := bp.buffer.Release(ctx, events); releaseErr != nil {
				bp.logger.Error("failed to release events", "error", releaseErr)
			}

			return fmt.Errorf("failed to write partition %v: %w", key, err)
		}
	}

	if err := bp.buffer.Ack(ctx, events); err != nil {
		return fmt.Errorf("failed to ack events: %w", err)
	}

	bp.logger.Debug("batch processed successfully", "events", len(events), "partitions", len(partitions))

	return nil
}

func (bp *BatchProcessor) groupByPartition(events []buffer.EventEnvelope) map[partitionKey][]buffer.EventEnvelope {
	partitions := make(map[partitionKey][]buffer.EventEnvelope)

	for _, event := range events {
		key := newPartitionKey(event, bp.batchWindow)
		partitions[key] = append(partitions[key], event)
	}

	return partitions
}

func (bp *BatchProcessor) writePartition(key partitionKey, batchID string, events []buffer.EventEnvelope) error {
	path := catalog.BuildPath(
		key.ProjectKey,
		key.LaunchUUID,
		string(key.EntityType),
		key.EventDate,
		batchID,
	)

	bp.logger.Debug("writing partition", "path", path, "events", len(events))

	if err := bp.writer.WritePartition(key.EntityType, path, events); err != nil {
		return fmt.Errorf("failed to write partition: %w", err)
	}

	return nil
}

func (bp *BatchProcessor) generateBatchID(key partitionKey) string {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	seq := bp.seqMap[key] + 1
	bp.seqMap[key] = seq

	cutoff := time.Now().Add(-2 * time.Hour).Unix()

	for k := range bp.seqMap {
		if k.WindowStart < cutoff {
			delete(bp.seqMap, k)
		}
	}

	return fmt.Sprintf("%d-%03d", key.WindowStart, seq)
}
