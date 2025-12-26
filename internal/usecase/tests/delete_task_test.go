package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/solumD/tasks-service/internal/usecase"
	"github.com/solumD/tasks-service/internal/usecase/mock"
	"github.com/solumD/tasks-service/pkg/logger"
)

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name                 string
		id                   int
		existFunc            func(ctx context.Context, id int) (bool, error)
		deleteFunc           func(ctx context.Context, id int) error
		expectedErr          error
		expectedExistCalled  bool
		expectedDeleteCalled bool
	}{
		{
			name: "task not exist",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return false, nil
			},
			deleteFunc:           nil,
			expectedErr:          usecase.ErrTaskNotFound,
			expectedExistCalled:  true,
			expectedDeleteCalled: false,
		},
		{
			name: "repo error on exist check",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return false, errors.New("db exist error")
			},
			deleteFunc:           nil,
			expectedErr:          errors.New("db exist error"),
			expectedExistCalled:  true,
			expectedDeleteCalled: false,
		},
		{
			name: "repo error on delete",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			deleteFunc: func(ctx context.Context, id int) error {
				return errors.New("db delete error")
			},
			expectedErr:          errors.New("db delete error"),
			expectedExistCalled:  true,
			expectedDeleteCalled: true,
		},
		{
			name: "success",
			id:   1,
			existFunc: func(ctx context.Context, id int) (bool, error) {
				return true, nil
			},
			deleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedErr:          nil,
			expectedExistCalled:  true,
			expectedDeleteCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockTaskRepo{
				IsTaskExistByIDFunc: tt.existFunc,
				DeleteTaskFunc:      tt.deleteFunc,
			}

			log := logger.NewMockLogger()
			u := usecase.NewTaskUsecase(repo, log)

			err := u.DeleteTask(context.Background(), tt.id)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if repo.IsTaskExistByIDCalled != tt.expectedExistCalled {
				t.Fatalf("expected IsTaskExistByID called = %v, got %v", tt.expectedExistCalled, repo.IsTaskExistByIDCalled)
			}

			if repo.DeleteTaskCalled != tt.expectedDeleteCalled {
				t.Fatalf("expected DeleteTask called = %v, got %v", tt.expectedDeleteCalled, repo.DeleteTaskCalled)
			}
		})
	}
}
