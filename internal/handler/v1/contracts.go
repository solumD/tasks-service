package v1

import (
	"context"

	"github.com/solumD/tasks-service/internal/model"
)

// TaskUsecase интерфейс изкейса Task
type TaskUsecase interface {
	CreateTask(ctx context.Context, task *model.Task) (int, error)
	GetAllTasks(ctx context.Context) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, id int) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id int) error
}
