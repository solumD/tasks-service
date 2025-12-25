package dto

import "github.com/solumD/tasks-service/internal/model"

func FromCreateReqToTask(req CreateTaskReq) *model.Task {
	return &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}
}

func FromTaskToResp(task *model.Task) *GetTaskByIDResp {
	return &GetTaskByIDResp{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
	}
}

func FromUpdateReqToTask(req UpdateTaskReq) *model.Task {
	return &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}
}

func FromTasksListToResp(tasks []*model.Task) *GetAllTasksResp {
	list := make([]*TaskDTO, 0, len(tasks))

	for _, task := range tasks {
		list = append(list, &TaskDTO{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Done:        task.Done,
		})
	}

	return &GetAllTasksResp{
		Tasks: list,
	}
}
