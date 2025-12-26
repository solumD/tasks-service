package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/solumD/tasks-service/internal/model"
	"github.com/solumD/tasks-service/internal/usecase"
	"github.com/solumD/tasks-service/internal/usecase/mock"
	"github.com/solumD/tasks-service/pkg/logger"
)

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name           string
		repoFunc       func(ctx context.Context) ([]*model.Task, error)
		expected       []*model.Task
		expectedErr    error
		expectedCalled bool
	}{
		{
			name: "repo returns error",
			repoFunc: func(ctx context.Context) ([]*model.Task, error) {
				return nil, errors.New("db error")
			},
			expected:       nil,
			expectedErr:    errors.New("db error"),
			expectedCalled: true,
		},
		{
			name: "success with sorting",
			repoFunc: func(ctx context.Context) ([]*model.Task, error) {
				return []*model.Task{
					{ID: 2, Title: "B"},
					{ID: 1, Title: "A"},
				}, nil
			},
			expected: []*model.Task{
				{ID: 1, Title: "A"},
				{ID: 2, Title: "B"},
			},
			expectedErr:    nil,
			expectedCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTaskRepo{
				GetAllTasksFunc: tt.repoFunc,
			}

			log := logger.NewMockLogger()
			u := usecase.NewTaskUsecase(repo, log)

			tasks, err := u.GetAllTasks(context.Background())

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if len(tasks) != len(tt.expected) {
				t.Fatalf("expected %d tasks, got %d", len(tt.expected), len(tasks))
			}

			for i := range tasks {
				if tasks[i].ID != tt.expected[i].ID || tasks[i].Title != tt.expected[i].Title {
					t.Fatalf("expected task %v, got %v", tt.expected[i], tasks[i])
				}
			}

			if repo.GetAllTasksCalled != tt.expectedCalled {
				t.Fatalf("expected GetAllTasks called = %v, got %v", tt.expectedCalled, repo.GetAllTasksCalled)
			}
		})
	}
}
