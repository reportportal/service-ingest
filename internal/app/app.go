package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/data/parquet"
	"github.com/reportportal/service-ingest/internal/data/repository"
	"github.com/reportportal/service-ingest/internal/handler"
	"github.com/reportportal/service-ingest/internal/processor"
	"github.com/reportportal/service-ingest/internal/service"
)

type App struct {
	server    *http.Server
	processor *processor.BatchProcessor
	buffer    buffer.Buffer
	ctx       context.Context
	cancel    context.CancelFunc
}

func New(cfg *config.Config) (*App, error) {
	flushInterval, err := cfg.Batch.FlushIntervalDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid flush interval: %w", err)
	}

	batchWindow, err := cfg.Batch.BatchWindowDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid batch window: %w", err)
	}

	buf, err := buffer.NewBadgerBuffer(cfg.Storage.BufferPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}

	writer := parquet.NewWriter(cfg.Storage.CatalogPath, cfg.Storage.ParquetCompression)

	launchRepo := repository.NewLaunchRepository(buf)
	itemRepo := repository.NewItemRepository(buf)
	logRepo := repository.NewLogRepository(buf)

	launchService := service.NewLaunchService(launchRepo)
	itemService := service.NewItemService(itemRepo)
	logService := service.NewLogService(logRepo)

	handlers := handler.Handlers{
		Launch: handler.NewLaunchHandler(launchService),
		Item:   handler.NewItemHandler(itemService),
		Log:    handler.NewLogHandler(logService),
	}

	router := handler.NewRouter(cfg.Server.BasePath, handlers)

	server := &http.Server{
		Addr:    cfg.Server.Addr(),
		Handler: router,
	}

	batchProcessor := processor.NewBatchProcessor(processor.BatchProcessorOptions{
		Buffer:        buf,
		Writer:        writer,
		FlushInterval: flushInterval,
		BatchWindow:   batchWindow,
		ReadLimit:     cfg.Batch.ReadLimit,
		Logger:        slog.Default(),
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		server:    server,
		processor: batchProcessor,
		buffer:    buf,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (a *App) Run() error {
	go a.processor.Start(a.ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	errChan := make(chan error, 1)

	go func() {
		slog.Info("http server listening", "addr", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		a.Shutdown()
		return fmt.Errorf("server error: %w", err)
	case sig := <-signalChan:
		slog.Info("received shutdown signal", "signal", sig)
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	slog.Info("shutting down gracefully...")

	a.cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	select {
	case <-a.processor.Done():
	case <-time.After(10 * time.Second):
		slog.Warn("batch processor shutdown timeout")
	}

	if err := a.buffer.Close(); err != nil {
		slog.Error("buffer close error", "error", err)
		return err
	}

	slog.Info("shutdown complete")
	return nil
}
