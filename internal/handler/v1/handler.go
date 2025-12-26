package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/solumD/tasks-service/internal/handler/v1/dto"
)

const (
	contentTypeEmpty = ""
	contentTypeJSON  = "application/json"
)

var (
	ErrFailedToDecodeReq   = errors.New("failed to decode request")
	ErrFailedToCreateTask  = errors.New("failed to create task")
	ErrFailedToGetTaskByID = errors.New("failed to get task by id")
	ErrFailedToUpdateTask  = errors.New("failed to update task")
	ErrFailedToDeleteTask  = errors.New("failed to delete task")
	ErrFailedToGetAllTasks = errors.New("failed to get all tasks")
	ErrInvalidTaskIDType   = errors.New("invalid task id type")
)

type handler struct {
	taskUsecase TaskUsecase
	log         *slog.Logger
}

func NewHandler(taskUsecase TaskUsecase, log *slog.Logger) *handler {
	return &handler{
		taskUsecase: taskUsecase,
		log:         log,
	}
}

func (h *handler) response(w http.ResponseWriter, contentType string, statusCode int, body []byte) {
	if len(contentType) > 0 {
		w.Header().Add("Content-Type", contentType)
	}

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	w.Write(body)
}

func (h *handler) errorResponse(w http.ResponseWriter, contentType string, statusCode int, err error) {
	body, errMarsh := json.Marshal(dto.NewErrorResponse(err.Error()))
	if errMarsh != nil {
		h.response(w, contentTypeEmpty, http.StatusInternalServerError, nil)
		return
	}

	if len(contentType) > 0 {
		w.Header().Add("Content-Type", contentType)
	}

	if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	w.Write(body)
}
