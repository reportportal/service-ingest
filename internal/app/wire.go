package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet"
	"github.com/reportportal/service-ingest/internal/data/repository"
	"github.com/reportportal/service-ingest/internal/handler"
	"github.com/reportportal/service-ingest/internal/processor"
	"github.com/reportportal/service-ingest/internal/service"
	"github.com/reportportal/service-ingest/pkg/logger"
)

func New(cfg *config.Config) (*App, error) {
	buf, err := buildBuffer(cfg)
	if err != nil {
		return nil, err
	}

	fileBuf := buildFileBuffer(cfg)

	writer := buildWriter(cfg)

	batchProcessor, err := buildBatchProcessor(cfg, buf, writer)
	if err != nil {
		_ = buf.Close()
		return nil, err
	}

	fileProcessor, err := buildFileProcessor(cfg, fileBuf)
	if err != nil {
		_ = buf.Close()
		return nil, err
	}

	handlers := buildHandlers(buf, fileBuf)
	server := buildServer(cfg, handlers)

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		server:        server,
		processor:     batchProcessor,
		fileProcessor: fileProcessor,
		buffer:        buf,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func buildBuffer(cfg *config.Config) (buffer.Buffer, error) {
	var opts badger.Options
	path := cfg.Buffer.BufferPath

	if path == "" {
		opts = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opts = badger.DefaultOptions(path)
	}

	opts = opts.WithBlockCacheSize(cfg.Buffer.GetBufferCacheSize())

	buf, err := buffer.NewBadgerBuffer(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}
	return buf, nil
}

func buildFileBuffer(cfg *config.Config) buffer.FileBuffer {
	return buffer.NewFileBuffer(cfg.Buffer.FileBufferPath)
}

func buildWriter(cfg *config.Config) *parquet.Writer {
	return parquet.NewWriter(cfg.Storage.CatalogPath, cfg.Storage.ParquetCompression)
}

func buildHandlers(buf buffer.Buffer, staging buffer.FileBuffer) handler.Handlers {
	launchRepo := repository.NewLaunchRepository(buf)
	itemRepo := repository.NewItemRepository(buf)
	logRepo := repository.NewLogRepository(buf, staging)

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
	level := logger.ParseLevel(cfg.Log.HTTPLevel)
	router := handler.NewRouter(cfg.Server.BasePath, handlers, level, cfg.Log.AddRSBody)
	return &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}
}

func buildBatchProcessor(cfg *config.Config, buf buffer.Buffer, writer *parquet.Writer) (*processor.BatchProcessor, error) {
	flushInterval, err := cfg.Processor.FlushIntervalDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid flush interval: %w", err)
	}

	return processor.NewBatchProcessor(processor.BatchProcessorOptions{
		Buffer:        buf,
		Writer:        writer,
		FlushInterval: flushInterval,
		ReadLimit:     cfg.Processor.ReadLimit,
		Logger:        slog.Default(),
	}), nil
}

func buildFileProcessor(cfg *config.Config, buffer buffer.FileBuffer) (*processor.FileProcessor, error) {
	if filepath.Clean(cfg.Buffer.FileBufferPath) == filepath.Clean(cfg.Storage.CatalogPath) {
		slog.Info("file processor disabled: buffer equals catalog path")
		return nil, nil
	}

	interval, err := cfg.Processor.FilesFlushIntervalDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid flush interval: %w", err)
	}

	return processor.NewFileProcessor(buffer, cfg.Storage.CatalogPath, interval), nil
}
