package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/solumD/tasks-service/internal/config"
	hnd "github.com/solumD/tasks-service/internal/handler"
	v1 "github.com/solumD/tasks-service/internal/handler/v1"
	inmemory "github.com/solumD/tasks-service/internal/repository/in_memory"
	"github.com/solumD/tasks-service/internal/usecase"
	httpserver "github.com/solumD/tasks-service/pkg/http_server"
	"github.com/solumD/tasks-service/pkg/logger"
)

const (
	shutdownTimeout = 10 * time.Second
)

func InitAndRun(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.LoggerLevel())

	taskRepo := inmemory.NewTaskRepo()
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	handler := v1.NewHandler(taskUsecase)

	r := hnd.NewRouter(ctx, log, handler)

	server := httpserver.New(cfg.ServerAddr(), r)
	server.Run()
	log.Info("starting server", slog.String("server address", cfg.ServerAddr()))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	log.Info("shutting down application")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdownCtx()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Error("error while shutting down server", logger.Err(err))
	}
}
