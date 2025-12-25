package usecase

import (
	"context"

	"github.com/solumD/tasks-service/internal/model"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, task *model.Task) (int, error)
	GetTaskByID(ctx context.Context, id int) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id int) error
	GetAllTasks(ctx context.Context) ([]*model.Task, error)
	IsTaskExistByID(ctx context.Context, id int) (bool, error)
}
