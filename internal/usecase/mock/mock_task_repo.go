package mock

import (
	"context"

	"github.com/solumD/tasks-service/internal/model"
)

// MockTaskRepo мок репозитория Task
type MockTaskRepo struct {
	CreateTaskFunc   func(ctx context.Context, task *model.Task) (int, error)
	CreateTaskCalled bool
	CreateTaskTask   *model.Task

	GetAllTasksFunc   func(ctx context.Context) ([]*model.Task, error)
	GetAllTasksCalled bool

	GetTaskByIDFunc   func(ctx context.Context, id int) (*model.Task, error)
	GetTaskByIDCalled bool
	GetTaskByIDID     int

	UpdateTaskFunc   func(ctx context.Context, task *model.Task) error
	UpdateTaskCalled bool
	UpdateTaskTask   *model.Task

	DeleteTaskFunc   func(ctx context.Context, id int) error
	DeleteTaskCalled bool
	DeleteTaskID     int

	IsTaskExistByIDFunc   func(ctx context.Context, id int) (bool, error)
	IsTaskExistByIDCalled bool
	IsTaskExistByIDID     int
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, task *model.Task) (int, error) {
	m.CreateTaskCalled = true
	m.CreateTaskTask = task

	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(ctx, task)
	}

	return 0, nil
}

func (m *MockTaskRepo) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	m.GetAllTasksCalled = true

	if m.GetAllTasksFunc != nil {
		return m.GetAllTasksFunc(ctx)
	}

	return nil, nil
}

func (m *MockTaskRepo) GetTaskByID(ctx context.Context, id int) (*model.Task, error) {
	m.GetTaskByIDCalled = true
	m.GetTaskByIDID = id

	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, id)
	}

	return nil, nil
}

func (m *MockTaskRepo) UpdateTask(ctx context.Context, task *model.Task) error {
	m.UpdateTaskCalled = true
	m.UpdateTaskTask = task

	if m.UpdateTaskFunc != nil {
		return m.UpdateTaskFunc(ctx, task)
	}

	return nil
}

func (m *MockTaskRepo) DeleteTask(ctx context.Context, id int) error {
	m.DeleteTaskCalled = true
	m.DeleteTaskID = id

	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}

	return nil
}

func (m *MockTaskRepo) IsTaskExistByID(ctx context.Context, id int) (bool, error) {
	m.IsTaskExistByIDCalled = true
	m.IsTaskExistByIDID = id

	if m.IsTaskExistByIDFunc != nil {
		return m.IsTaskExistByIDFunc(ctx, id)
	}

	return false, nil
}
