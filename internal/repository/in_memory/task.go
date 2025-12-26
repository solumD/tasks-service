package inmemory

import (
	"context"
	"sync"

	"github.com/solumD/tasks-service/internal/model"
)

type taskRepo struct {
	tasks map[int]*model.Task

	mu        *sync.RWMutex
	idCounter int
}

func NewTaskRepo() *taskRepo {
	return &taskRepo{
		tasks:     make(map[int]*model.Task),
		mu:        &sync.RWMutex{},
		idCounter: 0,
	}
}

// CreateTask создает новую задачу в хранилище
func (r *taskRepo) CreateTask(_ context.Context, task *model.Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.idCounter++
	task.ID = r.idCounter
	r.tasks[task.ID] = task

	return task.ID, nil
}

// GetAllTasks возвращает все задачи из хранилища
func (r *taskRepo) GetAllTasks(_ context.Context) ([]*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(r.tasks))

	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID возвращает задачу по ID из хранилища
func (r *taskRepo) GetTaskByID(_ context.Context, id int) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task := r.tasks[id]

	return task, nil
}

// UpdateTask обновляет задачу в хранилище
func (r taskRepo) UpdateTask(_ context.Context, task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldTask := r.tasks[task.ID]
	if len(task.Title) > 0 {
		oldTask.Title = task.Title
	}

	if len(task.Description) > 0 {
		oldTask.Description = task.Description
	}

	if task.Done {
		oldTask.Done = task.Done
	}

	r.tasks[oldTask.ID] = oldTask

	return nil
}

// DeleteTask удаляет задачу из хранилища
func (r *taskRepo) DeleteTask(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)

	return nil
}

func (r *taskRepo) IsTaskExistByID(_ context.Context, id int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exist := r.tasks[id]

	return exist, nil
}
