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

func TestGetTaskByID(t *testing.T) {
	tests := []struct {
		name                  string
		id                    int
		existFunc             func(ctx context.Context, id int) (bool, error)
		getByIDFunc           func(ctx context.Context, id int) (*model.Task, error)
		expected              *model.Task
		expectedErr           error
		expectedExistCalled   bool
		expectedGetByIDCalled bool
	}{
		{
			name: "task not exist",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return false, nil
			},
			getByIDFunc:           nil,
			expected:              nil,
			expectedErr:           usecase.ErrTaskNotFound,
			expectedExistCalled:   true,
			expectedGetByIDCalled: false,
		},
		{
			name: "repo error on exist check",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return false, errors.New("db error")
			},
			getByIDFunc:           nil,
			expected:              nil,
			expectedErr:           errors.New("db error"),
			expectedExistCalled:   true,
			expectedGetByIDCalled: false,
		},
		{
			name: "repo error on get task",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			getByIDFunc: func(ctx context.Context, id int) (*model.Task, error) {
				return nil, errors.New("db get error")
			},
			expected:              nil,
			expectedErr:           errors.New("db get error"),
			expectedExistCalled:   true,
			expectedGetByIDCalled: true,
		},
		{
			name: "success",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			getByIDFunc: func(ctx context.Context, id int) (*model.Task, error) {
				return &model.Task{ID: 1, Title: "Task1"}, nil
			},
			expected:              &model.Task{ID: 1, Title: "Task1"},
			expectedErr:           nil,
			expectedExistCalled:   true,
			expectedGetByIDCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTaskRepo{
				IsTaskExistByIDFunc: tt.existFunc,
				GetTaskByIDFunc:     tt.getByIDFunc,
			}

			log := logger.NewMockLogger()
			u := usecase.NewTaskUsecase(repo, log)

			task, err := u.GetTaskByID(context.Background(), tt.id)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expected != nil && (task.ID != tt.expected.ID || task.Title != tt.expected.Title) {
				t.Fatalf("expected task %v, got %v", tt.expected, task)
			}

			if repo.IsTaskExistByIDCalled != tt.expectedExistCalled {
				t.Fatalf("expected IsTaskExistByID called = %v, got %v", tt.expectedExistCalled, repo.IsTaskExistByIDCalled)
			}

			if repo.GetTaskByIDCalled != tt.expectedGetByIDCalled {
				t.Fatalf("expected GetTaskByID called = %v, got %v", tt.expectedGetByIDCalled, repo.GetTaskByIDCalled)
			}
		})
	}
}
