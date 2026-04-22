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

	opendal "github.com/apache/opendal/bindings/go"
	"github.com/reportportal/service-ingest/internal/data/buffer"
	"github.com/reportportal/service-ingest/internal/processor"
)

type App struct {
	server        *http.Server
	processor     *processor.BatchProcessor
	buffer        buffer.Buffer
	ctx           context.Context
	cancel        context.CancelFunc
	fileProcessor *processor.FileProcessor
	operator      *opendal.Operator
}

func (a *App) Run() error {
	go a.processor.Start(a.ctx)
	if a.fileProcessor != nil {
		go a.fileProcessor.Start(a.ctx)
	}

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
		if err = a.Shutdown(); err != nil {
			return err
		}
		return fmt.Errorf("server error: %w", err)
	case sig := <-signalChan:
		slog.Info("received shutdown signal", "signal", sig)
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	slog.Info("shutting down gracefully...")

	a.cancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	select {
	case <-a.processor.Done():
	case <-shutdownCtx.Done():
		slog.Warn("batch processor shutdown timeout")
	}

	if a.fileProcessor != nil {
		select {
		case <-a.fileProcessor.Done():
		case <-shutdownCtx.Done():
			slog.Warn("file processor shutdown timeout")
		}
	}

	a.operator.Close()

	if err := a.buffer.Close(); err != nil {
		slog.Error("buffer close error", "error", err)
		return err
	}

	slog.Info("shutdown complete")
	return nil
}
