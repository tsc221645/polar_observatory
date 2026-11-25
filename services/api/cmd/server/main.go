package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tsc221645/polar_observatory/services/api/internal/transport"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("POLAR_API_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8000"
	}
	if err := run(ctx, addr, os.Getenv("POLAR_API_DATABASE_URL"), logger); err != nil {
		logger.Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, addr, databaseURL string, logger *slog.Logger) error {
	var checkReady func(context.Context) bool
	if databaseURL != "" {
		config, err := pgxpool.ParseConfig(databaseURL)
		if err != nil {
			return errors.New("invalid POLAR_API_DATABASE_URL or pool configuration")
		}
		config.MaxConns = 4
		config.MinConns = 0
		config.MinIdleConns = 0
		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			return errors.New("could not initialize PostgreSQL pool")
		}
		defer pool.Close()
		checkReady = func(ctx context.Context) bool {
			// Query the migrated schema through the same role used for API reads.
			_, err := pool.Exec(ctx, "SELECT slug FROM polar.sources LIMIT 0")
			return err == nil
		}
	}
	return serve(ctx, addr, logger, transport.NewHandler(checkReady))
}

func serve(ctx context.Context, addr string, logger *slog.Logger, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       time.Minute,
	}
	done := make(chan error, 1)
	go func() {
		done <- server.Serve(listener)
	}()
	logger.Info("API listening", "address", listener.Addr().String())

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
