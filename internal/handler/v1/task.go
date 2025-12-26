package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/solumD/tasks-service/internal/handler/v1/dto"
	"github.com/solumD/tasks-service/internal/usecase"
	"github.com/solumD/tasks-service/pkg/logger"
)

// CreateTask обрабатывает запрос на создание новой задачи
func (h *handler) CreateTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.CreateTask"
		log := h.log.With(logger.String("fn", fn))

		log.Info("new request")

		var req dto.CreateTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrFailedToDecodeReq)
			return
		}

		log.Info("decoded request", logger.Any("request body", req))

		id, err := h.taskUsecase.CreateTask(ctx, dto.FromCreateReqToTask(req))
		if err != nil {
			if errors.Is(err, usecase.ErrEmptyTitle) {
				log.Error("failed to create task", logger.Error(err))

				h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, err)
				return
			}

			log.Error("failed to create task", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToCreateTask)
			return
		}

		resp := &dto.CreateTaskResp{ID: id}
		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Error("failed to marshal response", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToCreateTask)
			return
		}

		log.Info("created task", logger.Int("task id", id))

		h.response(w, contentTypeJSON, http.StatusCreated, respBody)
	}
}

// GetAllTasks обрабатывает запрос на получение всех задач
func (h *handler) GetAllTasks(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.GetAllTasks"
		log := h.log.With(logger.String("fn", fn))

		log.Info("new request")

		tasks, err := h.taskUsecase.GetAllTasks(ctx)
		if err != nil {
			log.Error("failed to get all tasks", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetAllTasks)
			return
		}

		resp := dto.FromTasksListToResp(tasks)
		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Error("failed to marshal response", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetAllTasks)
			return
		}

		log.Info("got all tasks", logger.Int("tasks count", len(tasks)))

		h.response(w, contentTypeJSON, http.StatusOK, respBody)
	}
}

// GetTaskByID обрабатывает запрос на получение задачи по id
func (h *handler) GetTaskByID(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.GetTaskByID"
		log := h.log.With(logger.String("fn", fn))

		log.Info("new request")

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		taskIDStr := pathParts[len(pathParts)-1]

		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Error("failed to get task id from path", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
			return
		}

		log.Info("got task id from path", logger.Int("task id", taskID))

		task, err := h.taskUsecase.GetTaskByID(ctx, taskID)
		if err != nil {
			if errors.Is(err, usecase.ErrTaskNotFound) {
				log.Error("failed to get task by id", logger.Error(err))

				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)
				return
			}

			log.Error("failed to get task by id", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetTaskByID)
			return
		}

		resp := dto.FromTaskToResp(task)
		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Error("failed to marshal response", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetTaskByID)
			return
		}

		log.Info("got task by id", logger.Int("task id", task.ID))

		h.response(w, contentTypeJSON, http.StatusOK, respBody)
	}
}

// UpdateTask обрабатывает запрос на обновление задачи
func (h *handler) UpdateTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.UpdateTask"
		log := h.log.With(logger.String("fn", fn))

		log.Info("new request")

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		taskIDStr := pathParts[len(pathParts)-1]

		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Error("failed to get task id from path", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
			return
		}

		log.Info("got task id from path", logger.Int("task id", taskID))

		var req dto.UpdateTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrFailedToDecodeReq)
			return
		}

		log.Info("decoded request", logger.Any("request body", req))

		task := dto.FromUpdateReqToTask(req)
		task.ID = taskID

		err = h.taskUsecase.UpdateTask(ctx, task)
		if err != nil {
			if errors.Is(err, usecase.ErrEmptyTitle) {
				log.Error("failed to update task", logger.Error(err))

				h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, err)
				return
			}

			if errors.Is(err, usecase.ErrTaskNotFound) {
				log.Error("failed to update task", logger.Error(err))

				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)
				return
			}

			log.Error("failed to update task", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToUpdateTask)
			return
		}

		log.Info("updated task", logger.Int("task id", task.ID))

		h.response(w, contentTypeJSON, http.StatusOK, nil)
	}
}

// DeleteTask обрабатывает запрос на удаление задачи
func (h *handler) DeleteTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handler.DeleteTask"
		log := h.log.With(logger.String("fn", fn))

		log.Info("new request")

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		taskIDStr := pathParts[len(pathParts)-1]

		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Error("failed to get task id from path", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
			return
		}

		log.Info("got task id from path", logger.Int("task id", taskID))

		err = h.taskUsecase.DeleteTask(ctx, taskID)
		if err != nil {
			if errors.Is(err, usecase.ErrTaskNotFound) {
				log.Error("failed to delete task", logger.Error(err))

				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)
				return
			}

			log.Error("failed to delete task", logger.Error(err))

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToDeleteTask)
			return
		}

		log.Info("deleted task", logger.Int("task id", taskID))

		h.response(w, contentTypeJSON, http.StatusOK, nil)
	}
}
