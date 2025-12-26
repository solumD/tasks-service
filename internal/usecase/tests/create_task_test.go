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

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name        string
		task        *model.Task
		repoFunc    func(ctx context.Context, task *model.Task) (int, error)
		expectedID  int
		expectedErr error
		repoCalled  bool
	}{
		{
			name: "empty title",
			task: &model.Task{Title: ""},
			repoFunc: nil,
			expectedID:  0,
			expectedErr: usecase.ErrEmptyTitle,
			repoCalled:  false,
		},
		{
			name: "repo returns error",
			task: &model.Task{Title: "task1"},
			repoFunc: func(ctx context.Context, task *model.Task) (int, error) {
				return 0, errors.New("db error")
			},
			expectedID:  0,
			expectedErr: errors.New("db error"),
			repoCalled:  true,
		},
		{
			name: "success",
			task: &model.Task{Title: "task2"},
			repoFunc: func(ctx context.Context, task *model.Task) (int, error) {
				return 42, nil
			},
			expectedID:  42,
			expectedErr: nil,
			repoCalled:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTaskRepo{
				CreateTaskFunc: tt.repoFunc,
			}

			log := logger.NewMockLogger()
			u := usecase.NewTaskUsecase(repo, log)

			id, err := u.CreateTask(context.Background(), tt.task)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if id != tt.expectedID {
				t.Fatalf("expected id %d, got %d", tt.expectedID, id)
			}

			if repo.CreateTaskCalled != tt.repoCalled {
				t.Fatalf("expected repo called = %v, got %v", tt.repoCalled, repo.CreateTaskCalled)
			}
		})
	}
}
