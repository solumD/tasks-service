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
	"github.com/solumD/tasks-service/internal/usecase"
	"github.com/solumD/tasks-service/pkg/logger"
)

func TestHandler_DeleteTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		pathID               string
		usecaseFunc          func(ctx context.Context, id int) error
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
			usecaseFunc: func(ctx context.Context, id int) error {
				return errors.New("db error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRespContains: "failed to delete task",
			expectedCalled:       true,
		},
		{
			name:   "not found",
			pathID: "1",
			usecaseFunc: func(ctx context.Context, id int) error {
				return usecase.ErrTaskNotFound
			},
			expectedStatus:       http.StatusNotFound,
			expectedRespContains: "task not found",
			expectedCalled:       true,
		},
		{
			name:   "success",
			pathID: "1",
			usecaseFunc: func(ctx context.Context, id int) error {
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
				DeleteTaskFunc: tt.usecaseFunc,
			}

			log := logger.NewMockLogger()
			h := v1.NewHandler(mockUsecase, log)

			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.pathID, nil)
			w := httptest.NewRecorder()

			h.DeleteTask(ctx).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedRespContains != "" && !strings.Contains(w.Body.String(), tt.expectedRespContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedRespContains, w.Body.String())
			}

			if mockUsecase.DeleteTaskCalled != tt.expectedCalled {
				t.Fatalf("expected DeleteTask called = %v, got %v", tt.expectedCalled, mockUsecase.DeleteTaskCalled)
			}
		})
	}
}
