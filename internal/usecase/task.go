package usecase

import (
	"context"
	"errors"
	"log/slog"
	"sort"

	"github.com/solumD/tasks-service/internal/model"
	"github.com/solumD/tasks-service/pkg/logger"
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
	const fn = "taskUsecase.CreateTask"
	log := u.log.With(logger.String("fn", fn))

	if len(task.Title) == 0 {
		return 0, ErrEmptyTitle
	}

	id, err := u.taskRepo.CreateTask(ctx, task)
	if err != nil {
		log.Error("failed to create task in repo", logger.Error(err))

		return 0, err
	}

	log.Info("created task in repo", logger.Int("task id", id))

	return id, nil
}

func (u *taskUsecase) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	const fn = "taskUsecase.GetAllTasks"
	log := u.log.With(logger.String("fn", fn))

	tasks, err := u.taskRepo.GetAllTasks(ctx)
	if err != nil {
		log.Error("failed to get all tasks from repo", logger.Error(err))

		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	log.Info("got all tasks from repo", logger.Int("tasks count", len(tasks)))

	return tasks, nil
}

func (u *taskUsecase) GetTaskByID(ctx context.Context, id int) (*model.Task, error) {
	const fn = "taskUsecase.GetTaskByID"
	log := u.log.With(logger.String("fn", fn))

	exist, err := u.taskRepo.IsTaskExistByID(ctx, id)
	if err != nil {
		log.Error("failed to check if task exist in repo", logger.Error(err))

		return nil, err
	}

	if !exist {
		return nil, ErrTaskNotFound
	}

	task, err := u.taskRepo.GetTaskByID(ctx, id)
	if err != nil {
		log.Error("failed to get task from repo", logger.Error(err))

		return nil, err
	}

	log.Info("got task from repo", logger.Int("task id", id))

	return task, nil
}

func (u *taskUsecase) UpdateTask(ctx context.Context, task *model.Task) error {
	const fn = "taskUsecase.UpdateTask"
	log := u.log.With(logger.String("fn", fn))

	exist, err := u.taskRepo.IsTaskExistByID(ctx, task.ID)
	if err != nil {
		log.Error("failed to check if task exist in repo", logger.Error(err))

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
		log.Error("failed to update task in repo", logger.Error(err))

		return err
	}

	log.Info("updated task in repo", logger.Int("task id", task.ID))

	return nil
}

func (u *taskUsecase) DeleteTask(ctx context.Context, id int) error {
	const fn = "taskUsecase.DeleteTask"
	log := u.log.With(logger.String("fn", fn))

	exist, err := u.taskRepo.IsTaskExistByID(ctx, id)
	if err != nil {
		log.Error("failed to check if task exist in repo", logger.Error(err))

		return err
	}

	if !exist {
		return ErrTaskNotFound
	}

	err = u.taskRepo.DeleteTask(ctx, id)
	if err != nil {
		log.Error("failed to delete task in repo", logger.Error(err))

		return err
	}

	log.Info("deleted task in repo", logger.Int("task id", id))

	return nil
}
