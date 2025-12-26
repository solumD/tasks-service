package handler

import (
	"context"
	"net/http"
)

// Handler интерфейс обработчика
type Handler interface {
	CreateTask(ctx context.Context) http.HandlerFunc
	GetAllTasks(ctx context.Context) http.HandlerFunc
	GetTaskByID(ctx context.Context) http.HandlerFunc
	UpdateTask(ctx context.Context) http.HandlerFunc
	DeleteTask(ctx context.Context) http.HandlerFunc
}
