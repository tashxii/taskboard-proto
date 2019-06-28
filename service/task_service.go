package service

import (
	"taskboard/model"
	"taskboard/orm"
	"taskboard/repository"

	"github.com/jinzhu/gorm"
)

// TaskService provides apis for task management.
type TaskService struct {
	tx            *gorm.DB
	taskRepo      *repository.TaskRepository
	taskEntryRepo *repository.TaskEntryRepository
}

// NewTaskService return new instance of TaskService.
func NewTaskService(tx *gorm.DB) *TaskService {
	return &TaskService{
		tx:            tx,
		taskRepo:      repository.NewTaskRepository(tx),
		taskEntryRepo: repository.NewTaskEntryRepository(tx),
	}
}

// FindTask returns task matching specified condition
func (s *TaskService) FindTask(condition interface{}) (*model.Task, error) {
	find, err := s.taskRepo.FindFirstTask(condition, []string{"id"})
	if err != nil {
		if err == orm.ErrorRecordNotFound {
			return nil, NewSvcErrorf(ErrorCodeNotFound, err, "Task not found")
		}
		return nil, NewSvcError(ErrorCodeDB, err, "Failed to find task")
	}
	return &find, nil
}

// FindTasks finds all tasks
func (s *TaskService) FindTasks(sortOrders []string) ([]model.Task, error) {
	tasks, err := s.taskRepo.FindTasks(&model.Task{}, 0, orm.NoLimit, sortOrders)
	if err != nil {
		return nil, NewSvcError(ErrorCodeDB, err, "Failed to find tasks")
	}
	return tasks, nil
}

// CreateTask creates new task
func (s *TaskService) CreateTask(task *model.Task) error {
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return NewSvcError(ErrorCodeDB, err, "Failed to create task")
	}
	return nil
}

// UpdateTask updates specifed task
func (s *TaskService) UpdateTask(task *model.Task) error {
	err := s.taskRepo.UpdateTask(task)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to update task. ID:%s", task.ID)
	}
	return nil
}

// DeleteTask deletes specifed task
func (s *TaskService) DeleteTask(task *model.Task) error {
	err := s.taskRepo.DeleteTask(task)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to delete task. ID:%s", task.ID)
	}
	err = s.taskEntryRepo.DeleteTaskEntriesByTaskID(task.ID)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to delete task's entries. TaskID:%s", task.ID)
	}
	return nil
}
