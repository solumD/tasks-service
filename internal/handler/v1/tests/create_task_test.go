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

func TestCreateTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		reqBody              string
		usecaseFunc          func(ctx context.Context, task *model.Task) (int, error)
		expectedStatus       int
		expectedRespContains string
		expectedCalled bool
	}{
		{
			name:                 "invalid JSON",
			reqBody:              "{invalid",
			usecaseFunc:          nil,
			expectedStatus:       http.StatusBadRequest,
			expectedRespContains: "failed to decode request",
			expectedCalled: false,
		},
		{
			name:    "empty title",
			reqBody: `{"title":""}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) (int, error) {
				return 0, usecase.ErrEmptyTitle
			},
			expectedStatus:       http.StatusBadRequest,
			expectedRespContains: "task title is empty",
			expectedCalled: true,
		},
		{
			name:    "repo error",
			reqBody: `{"title":"task1"}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) (int, error) {
				return 0, errors.New("repo error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRespContains: "failed to create task",
			expectedCalled: true,
		},

		{
			name:    "success",
			reqBody: `{"title":"task1"}`,
			usecaseFunc: func(ctx context.Context, task *model.Task) (int, error) {
				return 1, nil
			},
			expectedStatus:       http.StatusCreated,
			expectedRespContains: `"id":1`,
			expectedCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mock.MockTaskUsecase{
				CreateTaskFunc: tt.usecaseFunc,
			}

			log := logger.NewMockLogger()
			h := v1.NewHandler(mockUsecase, log)

			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.reqBody))
			w := httptest.NewRecorder()

			h.CreateTask(ctx).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedRespContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedRespContains, w.Body.String())
			}

			if mockUsecase.CreateTaskCalled != tt.expectedCalled {
				t.Fatalf("expected CreateTask called = %v, got %v", tt.expectedCalled, mockUsecase.CreateTaskCalled)
			}
		})
	}
}
