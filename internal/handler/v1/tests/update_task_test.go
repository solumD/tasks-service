package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v1 "github.com/solumD/tasks-service/internal/handler/v1"
	"github.com/solumD/tasks-service/internal/handler/v1/mock"
	"github.com/solumD/tasks-service/internal/model"
	"github.com/solumD/tasks-service/internal/usecase"
	"github.com/solumD/tasks-service/pkg/logger"
)

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		pathID               string
		reqBody              string
		usecaseFunc          func(ctx context.Context, task *model.Task) error
		expectedStatus       int
		expectedRespContains string
		expectedCalled       bool
	}{
		{
			name:                 "invalid ID",
			pathID:               "abc",
			reqBody:              `{"title":"A"}`,
			usecaseFunc:          nil,
			expectedStatus:       http.StatusBadRequest,
			expectedRespContains: "invalid task id type",
			expectedCalled:       false,
		},
		{
			name:    "repo error",
			pathID:  "1",
			reqBody: `{"title":"Updated"}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) error {
				return errors.New("db error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRespContains: "failed to update task",
			expectedCalled:       true,
		},
		{
			name:    "not found",
			pathID:  "1",
			reqBody: `{"title":"Updated"}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) error {
				return usecase.ErrTaskNotFound
			},
			expectedStatus:       http.StatusNotFound,
			expectedRespContains: "task not found",
			expectedCalled:       true,
		},
		{
			name:    "empty title",
			pathID:  "1",
			reqBody: `{"title":""}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) error {
				return usecase.ErrEmptyTitle
			},
			expectedStatus:       http.StatusBadRequest,
			expectedRespContains: "task title is empty",
			expectedCalled:       true,
		},
		{
			name:    "success",
			pathID:  "1",
			reqBody: `{"title":"Updated"}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) error {
				return nil
			},
			expectedStatus:       http.StatusOK,
			expectedRespContains: "",
			expectedCalled:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mock.MockTaskUsecase{
				UpdateTaskFunc: tt.usecaseFunc,
			}

			log := logger.NewMockLogger()
			h := v1.NewHandler(mockUsecase, log)

			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.pathID, strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			h.UpdateTask(ctx).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedRespContains != "" && !strings.Contains(w.Body.String(), tt.expectedRespContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedRespContains, w.Body.String())
			}

			if mockUsecase.UpdateTaskCalled != tt.expectedCalled {
				t.Fatalf("expected UpdateTask called = %v, got %v", tt.expectedCalled, mockUsecase.UpdateTaskCalled)
			}
		})
	}
}
