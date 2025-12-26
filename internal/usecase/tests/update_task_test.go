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

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name                 string
		task                 *model.Task
		existFunc            func(ctx context.Context, id int) (bool, error)
		updateFunc           func(ctx context.Context, task *model.Task) error
		expectedErr          error
		expectedExistCalled  bool
		expectedUpdateCalled bool
	}{
		{
			name: "task not exist",
			task: &model.Task{ID: 1, Title: "Task1"},
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return false, nil
			},
			updateFunc:           nil,
			expectedErr:          usecase.ErrTaskNotFound,
			expectedExistCalled:  true,
			expectedUpdateCalled: false,
		},
		{
			name: "empty title",
			task: &model.Task{ID: 1, Title: ""},
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			updateFunc:           nil,
			expectedErr:          usecase.ErrEmptyTitle,
			expectedExistCalled:  true,
			expectedUpdateCalled: false,
		},
		{
			name: "repo error on update",
			task: &model.Task{ID: 1, Title: "Task1"},
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			updateFunc: func(ctx context.Context, task *model.Task) error {
				return errors.New("db update error")
			},
			expectedErr:          errors.New("db update error"),
			expectedExistCalled:  true,
			expectedUpdateCalled: true,
		},
		{
			name: "success",
			task: &model.Task{ID: 1, Title: "Task1"},
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			updateFunc: func(ctx context.Context, task *model.Task) error {
				return nil
			},
			expectedErr:          nil,
			expectedExistCalled:  true,
			expectedUpdateCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTaskRepo{
				IsTaskExistByIDFunc: tt.existFunc,
				UpdateTaskFunc:      tt.updateFunc,
			}

			log := logger.NewMockLogger()
			u := usecase.NewTaskUsecase(repo, log)
			err := u.UpdateTask(context.Background(), tt.task)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if repo.IsTaskExistByIDCalled != tt.expectedExistCalled {
				t.Fatalf("expected IsTaskExistByID called = %v, got %v", tt.expectedExistCalled, repo.IsTaskExistByIDCalled)
			}

			if repo.UpdateTaskCalled != tt.expectedUpdateCalled {
				t.Fatalf("expected UpdateTask called = %v, got %v", tt.expectedUpdateCalled, repo.UpdateTaskCalled)
			}
		})
	}
}
