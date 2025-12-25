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
)

func (h *handler) CreateTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTaskReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrFailedToDecodeReq)

			return
		}

		id, err := h.taskUsecase.CreateTask(ctx, dto.FromCreateReqToTask(req))
		if err != nil {
			if errors.Is(err, usecase.ErrEmptyTitle) {
				h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, err)

				return
			}

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToCreateTask)

			return
		}

		resp := &dto.CreateTaskResp{ID: id}
		respBody, err := json.Marshal(resp)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToCreateTask)

			return
		}

		h.response(w, contentTypeJSON, http.StatusCreated, respBody)
	}
}

func (h *handler) GetAllTasks(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.taskUsecase.GetAllTasks(ctx)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetAllTasks)

			return
		}

		resp := dto.FromTasksListToResp(tasks)
		respBody, err := json.Marshal(resp)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetAllTasks)

			return
		}

		h.response(w, contentTypeJSON, http.StatusOK, respBody)
	}
}

func (h *handler) GetTaskByID(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDstr := strings.TrimSpace(r.PathValue("id"))

		taskID, err := strconv.Atoi(taskIDstr)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
		}

		task, err := h.taskUsecase.GetTaskByID(ctx, taskID)
		if err != nil {
			if errors.Is(err, usecase.ErrTaskNotFound) {
				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)

				return
			}

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetTaskByID)
		}

		resp := dto.FromTaskToResp(task)
		respBody, err := json.Marshal(resp)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToGetTaskByID)

			return
		}

		h.response(w, contentTypeJSON, http.StatusOK, respBody)
	}
}

func (h *handler) UpdateTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDstr := strings.TrimSpace(r.PathValue("id"))

		taskID, err := strconv.Atoi(taskIDstr)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
		}

		var req dto.UpdateTaskReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrFailedToDecodeReq)

			return
		}

		task := dto.FromUpdateReqToTask(req)
		task.ID = taskID

		err = h.taskUsecase.UpdateTask(ctx, task)
		if err != nil {
			if errors.Is(err, usecase.ErrEmptyTitle) {
				h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, err)

				return
			}

			if errors.Is(err, usecase.ErrTaskNotFound) {
				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)

				return
			}

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToUpdateTask)

			return
		}

		h.response(w, contentTypeJSON, http.StatusOK, nil)
	}
}

func (h *handler) DeleteTask(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDstr := strings.TrimSpace(r.PathValue("id"))

		taskID, err := strconv.Atoi(taskIDstr)
		if err != nil {
			h.errorResponse(w, contentTypeJSON, http.StatusBadRequest, ErrInvalidTaskIDType)
		}

		err = h.taskUsecase.DeleteTask(ctx, taskID)
		if err != nil {
			if errors.Is(err, usecase.ErrTaskNotFound) {
				h.errorResponse(w, contentTypeJSON, http.StatusNotFound, err)

				return
			}

			h.errorResponse(w, contentTypeJSON, http.StatusInternalServerError, ErrFailedToDeleteTask)

			return
		}

		h.response(w, contentTypeJSON, http.StatusOK, nil)
	}
}
