package handler

import (
	"context"
	"net/http"
)

type Handler interface {
	CreateTask(ctx context.Context) http.HandlerFunc
	GetAllTasks(ctx context.Context) http.HandlerFunc
	GetTaskByID(ctx context.Context) http.HandlerFunc
	UpdateTask(ctx context.Context) http.HandlerFunc
	DeleteTask(ctx context.Context) http.HandlerFunc
}

func NewRouter(ctx context.Context, handler Handler) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("POST /tasks", handler.CreateTask(ctx))
	r.HandleFunc("GET /tasks", handler.GetAllTasks(ctx))
	r.HandleFunc("GET /tasks/{id}", handler.GetTaskByID(ctx))
	r.HandleFunc("PUT /tasks/{id}", handler.UpdateTask(ctx))
	r.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask(ctx))

	return r
}
