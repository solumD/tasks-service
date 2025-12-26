package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/solumD/tasks-service/pkg/middleware"
)

// NewRouter возвращает роутер для обработки запросов
func NewRouter(ctx context.Context, log *slog.Logger, handler Handler) *http.ServeMux {
	r := http.NewServeMux()
	loggerMW := middleware.NewMWLogger(log)

	r.Handle(
		"POST /todos",
		loggerMW(http.HandlerFunc(handler.CreateTask(ctx))),
	)

	r.Handle(
		"GET /todos",
		loggerMW(http.HandlerFunc(handler.GetAllTasks(ctx))),
	)

	r.Handle(
		"GET /todos/{id}",
		loggerMW(http.HandlerFunc(handler.GetTaskByID(ctx))),
	)

	r.Handle(
		"PUT /todos/{id}",
		loggerMW(http.HandlerFunc(handler.UpdateTask(ctx))),
	)

	r.Handle(
		"DELETE /todos/{id}",
		loggerMW(http.HandlerFunc(handler.DeleteTask(ctx))),
	)

	return r
}
