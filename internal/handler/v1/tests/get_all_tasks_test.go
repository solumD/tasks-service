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
	"github.com/solumD/tasks-service/pkg/logger"
)

func TestGetAllTasks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name                 string
		usecaseFunc          func(ctx context.Context) ([]*model.Task, error)
		expectedStatus       int
		expectedRespContains string
		expectedCalled       bool
	}{
		{
			name: "repo error",
			usecaseFunc: func(ctx context.Context) ([]*model.Task, error) {
				return nil, errors.New("db error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRespContains: "failed to get all tasks",
			expectedCalled:       true,
		},
		{
			name: "success",
			usecaseFunc: func(ctx context.Context) ([]*model.Task, error) {
				return []*model.Task{
					{ID: 1, Title: "A"},
					{ID: 2, Title: "B"},
				}, nil
			},
			expectedStatus:       http.StatusOK,
			expectedRespContains: `"title":"A"`,
			expectedCalled:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := &mock.MockTaskUsecase{
				GetAllTasksFunc: tt.usecaseFunc,
			}

			log := logger.NewMockLogger()
			h := v1.NewHandler(mockUsecase, log)

			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			w := httptest.NewRecorder()

			h.GetAllTasks(ctx).ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedRespContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedRespContains, w.Body.String())
			}

			if mockUsecase.GetAllTasksCalled != tt.expectedCalled {
				t.Fatalf("expected GetAllTasks called = %v, got %v", tt.expectedCalled, mockUsecase.GetAllTasksCalled)
			}
		})
	}
}
