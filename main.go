package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"gitlab.com/hmajid2301/goroutinely/config"
	"gitlab.com/hmajid2301/goroutinely/store/db"
	transporthttp "gitlab.com/hmajid2301/goroutinely/transport/http"
)

//go:embed static
var staticFiles embed.FS

func main() {
	var exitCode int

	err := mainLogic()
	if err != nil {
		slog.Error("failed to start app", "error", err)
		exitCode = 1
	}
	defer func() { os.Exit(exitCode) }()
}

func mainLogic() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	pool, err := db.NewPool(ctx, conf.DB.URI)
	if err != nil {
		return fmt.Errorf("failed to setup database pool: %w", err)
	}
	defer pool.Close()

	database := db.NewDB(pool)

	server := transporthttp.NewServer(logger, staticFiles, conf.Server)

	go func() {
		logger.InfoContext(
			ctx,
			"starting server",
			slog.String("host", conf.Server.Host),
			slog.Int("port", conf.Server.Port),
		)
		if err := server.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to serve server", "error", err)
		}
	}()

	terminateHandler(ctx, cancel, logger, server, database)

	return nil
}

func terminateHandler(
	ctx context.Context,
	cancel context.CancelFunc,
	logger *slog.Logger,
	srv *transporthttp.Server,
	_ *db.DB,
) {
	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	stop()
	logger.InfoContext(ctx, "received shutdown signal, starting graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.ErrorContext(ctx, "unexpected error while shutting down server", slog.Any("error", err))
	} else {
		logger.InfoContext(ctx, "server shutdown completed successfully")
	}
}
