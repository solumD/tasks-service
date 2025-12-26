package usecase

import (
	"context"
	"errors"
	"log/slog"
	"sort"

	"github.com/solumD/tasks-service/internal/model"
)

var (
	ErrEmptyTitle   = errors.New("task title is empty")
	ErrTaskNotFound = errors.New("task not found")
)

type taskUsecase struct {
	taskRepo TaskRepo
	log      *slog.Logger
}

func NewTaskUsecase(taskRepo TaskRepo, log *slog.Logger) *taskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
		log:      log,
	}
}

func (u *taskUsecase) CreateTask(ctx context.Context, task *model.Task) (int, error) {
	if len(task.Title) == 0 {
		return 0, ErrEmptyTitle
	}

	id, err := u.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *taskUsecase) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	tasks, err := u.taskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

func (u *taskUsecase) GetTaskByID(ctx context.Context, id int) (*model.Task, error) {
	exist, err := u.taskRepo.IsTaskExistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrTaskNotFound
	}

	task, err := u.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (u *taskUsecase) UpdateTask(ctx context.Context, task *model.Task) error {
	exist, err := u.taskRepo.IsTaskExistByID(ctx, task.ID)
	if err != nil {
		return err
	}

	if !exist {
		return ErrTaskNotFound
	}

	if len(task.Title) == 0 {
		return ErrEmptyTitle
	}

	err = u.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id int) error {
	exist, err := u.taskRepo.IsTaskExistByID(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		return ErrTaskNotFound
	}

	err = u.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
