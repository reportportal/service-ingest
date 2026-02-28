package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet"
	"github.com/reportportal/service-ingest/internal/data/repository"
	"github.com/reportportal/service-ingest/internal/handler"
	"github.com/reportportal/service-ingest/internal/processor"
	"github.com/reportportal/service-ingest/internal/service"
)

func New(cfg *config.Config) (*App, error) {
	buf, err := buildBuffer(cfg)
	if err != nil {
		return nil, err
	}

	writer := buildWriter(cfg)
	handlers := buildHandlers(buf)
	server := buildServer(cfg, handlers)

	batchProcessor, err := buildBatchProcessor(cfg, buf, writer)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		server:    server,
		processor: batchProcessor,
		buffer:    buf,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func buildBuffer(cfg *config.Config) (buffer.Buffer, error) {
	buf, err := buffer.NewBadgerBuffer(cfg.Storage.BufferPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}
	return buf, nil
}

func buildWriter(cfg *config.Config) *parquet.Writer {
	return parquet.NewWriter(cfg.Storage.CatalogPath, cfg.Storage.ParquetCompression)
}

func buildHandlers(buf buffer.Buffer) handler.Handlers {
	launchRepo := repository.NewLaunchRepository(buf)
	itemRepo := repository.NewItemRepository(buf)
	logRepo := repository.NewLogRepository(buf)

	launchService := service.NewLaunchService(launchRepo)
	itemService := service.NewItemService(itemRepo)
	logService := service.NewLogService(logRepo)

	return handler.Handlers{
		Launch: handler.NewLaunchHandler(launchService),
		Item:   handler.NewItemHandler(itemService),
		Log:    handler.NewLogHandler(logService),
	}
}

func buildServer(cfg *config.Config, handlers handler.Handlers) *http.Server {
	router := handler.NewRouter(cfg.Server.BasePath, handlers)
	return &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}
}

func buildBatchProcessor(cfg *config.Config, buf buffer.Buffer, writer *parquet.Writer) (*processor.BatchProcessor, error) {
	flushInterval, err := cfg.Batch.FlushIntervalDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid flush interval: %w", err)
	}

	batchWindow, err := cfg.Batch.BatchWindowDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid batch window: %w", err)
	}

	return processor.NewBatchProcessor(processor.BatchProcessorOptions{
		Buffer:        buf,
		Writer:        writer,
		FlushInterval: flushInterval,
		BatchWindow:   batchWindow,
		ReadLimit:     cfg.Batch.ReadLimit,
		Logger:        slog.Default(),
	}), nil
}
