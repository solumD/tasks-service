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

func TestGetTaskByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		pathID               string
		usecaseFunc          func(ctx context.Context, id int) (*model.Task, error)
		expectedStatus       int
		expectedRespContains string
		expectedCalled       bool
	}{
		{
			name:                 "invalid ID",
			pathID:               "abc",
			usecaseFunc:          nil,
			expectedStatus:       http.StatusBadRequest,
			expectedRespContains: "invalid task id type",
			expectedCalled:       false,
		},
		{
			name:   "repo error",
			pathID: "1",
			usecaseFunc: func(ctx context.Context, id int) (*model.Task, error) {
				return nil, errors.New("db error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRespContains: "failed to get task by id",
			expectedCalled:       true,
		},
		{
			name:   "not found",
			pathID: "1",
			usecaseFunc: func(ctx context.Context, id int) (*model.Task, error) {
				return nil, usecase.ErrTaskNotFound
			},
			expectedStatus:       http.StatusNotFound,
			expectedRespContains: "task not found",
			expectedCalled:       true,
		},
		{
			name:   "success",
			pathID: "1",
			usecaseFunc: func(ctx context.Context, id int) (*model.Task, error) {
				return &model.Task{ID: 1, Title: "Task1"}, nil
			},
			expectedStatus:       http.StatusOK,
			expectedRespContains: `"title":"Task1"`,
			expectedCalled:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mock.MockTaskUsecase{
				GetTaskByIDFunc: tt.usecaseFunc,
			}

			log := logger.NewMockLogger()
			h := v1.NewHandler(mockUsecase, log)

			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.pathID, nil)
			w := httptest.NewRecorder()

			h.GetTaskByID(ctx).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedRespContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedRespContains, w.Body.String())
			}

			if mockUsecase.GetTaskByIDCalled != tt.expectedCalled {
				t.Fatalf("expected GetTaskByID called = %v, got %v", tt.expectedCalled, mockUsecase.GetTaskByIDCalled)
			}

		})
	}
}
