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
		"POST /tasks",
		loggerMW(http.HandlerFunc(handler.CreateTask(ctx))),
	)

	r.Handle(
		"GET /tasks",
		loggerMW(http.HandlerFunc(handler.GetAllTasks(ctx))),
	)

	r.Handle(
		"GET /tasks/{id}",
		loggerMW(http.HandlerFunc(handler.GetTaskByID(ctx))),
	)

	r.Handle(
		"PUT /tasks/{id}",
		loggerMW(http.HandlerFunc(handler.UpdateTask(ctx))),
	)

	r.Handle(
		"DELETE /tasks/{id}",
		loggerMW(http.HandlerFunc(handler.DeleteTask(ctx))),
	)

	return r
}
