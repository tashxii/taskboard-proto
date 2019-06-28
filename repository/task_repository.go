package repository

import (
	"sync"
	"taskboard/model"
	"taskboard/orm"

	"github.com/jinzhu/gorm"
)

var lockTask = &sync.Mutex{}

// TaskRepository is repository of task table
type TaskRepository struct {
	tx *gorm.DB
}

// NewTaskRepository returns new instance of TaskRepository
func NewTaskRepository(tx *gorm.DB) *TaskRepository {
	if tx == nil {
		// Programing error!!
		panic("tx must be set")
	}
	return &TaskRepository{
		tx: tx,
	}
}

// FindFirstTask returns first Task matching with specified condition
func (repo *TaskRepository) FindFirstTask(condition interface{}, sortOrders []string) (result model.Task, err error) {
	query := repo.tx.Where(condition)
	if sortOrders == nil {
		sortOrders = []string{}
	}

	for _, sortOrder := range sortOrders {
		query = query.Order(sortOrder)
	}
	err = query.First(&result).Error
	return
}

// FindTasks returns Tasks matching with specified condition
func (repo *TaskRepository) FindTasks(condition interface{}, offset int, limit int, sortOrders []string) (result []model.Task, err error) {
	query := repo.tx.Where(condition)
	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit >= 0 {
		query = query.Limit(limit)
	}

	if sortOrders == nil {
		sortOrders = []string{}
	}
	for _, task := range sortOrders {
		query = query.Order(task)
	}

	err = query.Find(&result).Error
	return
}

// CountTasks returns the number of Tasks matching specfied condition
func (repo *TaskRepository) CountTasks(condition interface{}) (count int, err error) {
	var tasks []model.Task
	err = repo.tx.Where(condition).Find(&tasks).Count(&count).Error
	return
}

// CreateTask inserts new Task record
func (repo *TaskRepository) CreateTask(task *model.Task) error {
	return repo.CreateTasks([]*model.Task{task})
}

// UpdateTask updates Task record
func (repo *TaskRepository) UpdateTask(task *model.Task) error {
	return repo.UpdateTasks([]*model.Task{task})
}

// DeleteTask deletes Task record
func (repo *TaskRepository) DeleteTask(task *model.Task) error {
	return repo.DeleteTasks([]*model.Task{task})
}

// CreateTasks inserts new Task records.
func (repo *TaskRepository) CreateTasks(tasks []*model.Task) (err error) {
	for _, task := range tasks {
		err = repo.tx.Create(task).Error
		if err != nil {
			return
		}
	}
	return
}

// UpdateTasks updates task records
func (repo *TaskRepository) UpdateTasks(tasks []*model.Task) (err error) {
	lockTask.Lock()
	defer lockTask.Unlock()

	for _, task := range tasks {
		oldVersion := task.Version
		task.Version++
		db := repo.tx.Model(&model.Task{}).Where("version = ?", oldVersion).Updates(task)
		count := db.RowsAffected
		err = db.Error
		// return ErrorRecordNotFoud as optimistic lock error
		if err == nil && count == 0 {
			return orm.ErrorRecordNotFound
		}
		if err != nil {
			return
		}
	}
	return
}

// DeleteTasks deletes Task records
func (repo *TaskRepository) DeleteTasks(tasks []*model.Task) (err error) {
	for _, task := range tasks {
		if task.ID == "" {
			continue // To avoid deleting all due to gorm warning, continue here.
		}
		err = repo.tx.Delete(task).Error
		if err != nil {
			return
		}
	}
	return
}

// ClearBoardID clears board_id field which matches specified board id
func (repo *TaskRepository) ClearBoardID(boardID string) (err error) {
	return repo.tx.Table("tasks").Where("board_id = ?", boardID).UpdateColumn("board_id = ?", "").Error
}
