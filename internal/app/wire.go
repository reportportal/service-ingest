package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/apache/opendal-go-services/fs"
	"github.com/apache/opendal-go-services/s3"
	opendal "github.com/apache/opendal/bindings/go"
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
	op, err := newStorageOperator(cfg)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			op.Close()
		}
	}()

	buf, err := buildBuffer(cfg)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = buf.Close()
		}
	}()

	fileBuf := buildFileBuffer(cfg)

	batchProcessor, err := buildBatchProcessor(cfg, buf, op)
	if err != nil {
		return nil, err
	}

	fileProcessor, err := buildFileProcessor(cfg, fileBuf, op)
	if err != nil {
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
		operator:      op,
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

func buildBatchProcessor(cfg *config.Config, buf buffer.Buffer, operator *opendal.Operator) (*processor.BatchProcessor, error) {
	writer := parquet.NewWriter(cfg.Storage.CatalogPath, cfg.Storage.ParquetCompression, operator)
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

func buildFileProcessor(cfg *config.Config, buffer buffer.FileBuffer, operator *opendal.Operator) (*processor.FileProcessor, error) {
	if cfg.Storage.Type == "fs" && filepath.Clean(cfg.Buffer.FileBufferPath) == filepath.Clean(cfg.Storage.CatalogPath) {
		slog.Info("file processor disabled: buffer equals catalog path")
		return nil, nil
	}

	interval, err := cfg.Processor.FilesFlushIntervalDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid flush interval: %w", err)
	}

	return processor.NewFileProcessor(buffer, cfg.Storage.CatalogPath, interval, operator), nil
}

func newStorageOperator(cfg *config.Config) (*opendal.Operator, error) {
	switch cfg.Storage.Type {
	case "fs":
		opts := opendal.OperatorOptions{"root": cfg.Storage.CatalogPath}
		return opendal.NewOperator(fs.Scheme, opts)
	case "s3":
		return opendal.NewOperator(s3.Scheme, s3Options(cfg))
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}
}

func s3Options(cfg *config.Config) opendal.OperatorOptions {
	opts := opendal.OperatorOptions{
		"root":   cfg.Storage.CatalogPath,
		"bucket": cfg.S3.Bucket,
	}

	if cfg.S3.Region != "" {
		opts["region"] = cfg.S3.Region
	}

	if cfg.S3.Endpoint != "" {
		opts["endpoint"] = cfg.S3.Endpoint
	}

	if cfg.S3.AccessKey != "" {
		opts["access_key_id"] = cfg.S3.AccessKey
		opts["secret_access_key"] = cfg.S3.SecretKey
	}

	if cfg.S3.SessionToken != "" {
		opts["session_token"] = cfg.S3.SessionToken
	}

	return opts
}
