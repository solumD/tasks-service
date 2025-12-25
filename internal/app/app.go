package app

import (
	"context"
	"log"
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
)

const (
	shutdownTimeout = 10 * time.Second
)

func InitAndRun(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.MustLoad()

	taskRepo := inmemory.NewTaskRepo()
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	handler := v1.NewHandler(taskUsecase)

	r := hnd.NewRouter(ctx, handler)

	server := httpserver.New(cfg.ServerAddr(), r)
	server.Run()
	log.Printf("starting server on %s\n", cfg.ServerAddr())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	log.Print("shutting down app...\n")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdownCtx()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("error while shutting down server: %v\n", err)
	}
}
