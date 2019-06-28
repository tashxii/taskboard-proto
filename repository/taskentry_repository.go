package repository

import (
	"sync"
	"taskboard/model"

	"github.com/jinzhu/gorm"
)

var lockTaskEntry = &sync.Mutex{}

// TaskEntryRepository is repository of taskEntry table
type TaskEntryRepository struct {
	tx *gorm.DB
}

// NewTaskEntryRepository returns new instance of TaskEntryRepository
func NewTaskEntryRepository(tx *gorm.DB) *TaskEntryRepository {
	if tx == nil {
		// Programing error!!
		panic("tx must be set")
	}
	return &TaskEntryRepository{
		tx: tx,
	}
}

// FindFirstTaskEntry returns first TaskEntry matching with specified condition
func (repo *TaskEntryRepository) FindFirstTaskEntry(condition interface{}, sortOrders []string) (result model.TaskEntry, err error) {
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

// FindTaskEntries returns TaskEntries matching with specified condition
func (repo *TaskEntryRepository) FindTaskEntries(condition interface{}, offset int, limit int, sortOrders []string) (result []model.TaskEntry, err error) {
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
	for _, taskEntry := range sortOrders {
		query = query.Order(taskEntry)
	}

	err = query.Find(&result).Error
	return
}

// CountTaskEntries returns the number of TaskEntries matching specfied condition
func (repo *TaskEntryRepository) CountTaskEntries(condition interface{}) (count int, err error) {
	var taskEntries []model.TaskEntry
	err = repo.tx.Where(condition).Find(&taskEntries).Count(&count).Error
	return
}

// CreateTaskEntry inserts new TaskEntry record
func (repo *TaskEntryRepository) CreateTaskEntry(taskEntry *model.TaskEntry) error {
	return repo.CreateTaskEntries([]*model.TaskEntry{taskEntry})
}

// UpdateTaskEntry updates TaskEntry record
func (repo *TaskEntryRepository) UpdateTaskEntry(taskEntry *model.TaskEntry) error {
	return repo.UpdateTaskEntries([]*model.TaskEntry{taskEntry})
}

// DeleteTaskEntry deletes TaskEntry record
func (repo *TaskEntryRepository) DeleteTaskEntry(taskEntry *model.TaskEntry) error {
	return repo.DeleteTaskEntries([]*model.TaskEntry{taskEntry})
}

// CreateTaskEntries inserts new TaskEntry records.
func (repo *TaskEntryRepository) CreateTaskEntries(taskEntries []*model.TaskEntry) (err error) {
	for _, taskEntry := range taskEntries {
		err = repo.tx.Create(taskEntry).Error
		if err != nil {
			return
		}
	}
	return
}

// UpdateTaskEntries updates taskEntry records
func (repo *TaskEntryRepository) UpdateTaskEntries(taskEntries []*model.TaskEntry) (err error) {
	lockTaskEntry.Lock()
	defer lockTaskEntry.Unlock()

	for _, taskEntry := range taskEntries {
		err = repo.tx.Update(taskEntry).Error
		if err != nil {
			return
		}
	}
	return
}

// DeleteTaskEntries deletes TaskEntry records
func (repo *TaskEntryRepository) DeleteTaskEntries(taskEntries []*model.TaskEntry) (err error) {
	for _, taskEntry := range taskEntries {
		if taskEntry.ID == "" {
			continue // To avoid deleting all due to gorm warning, continue here.
		}
		err = repo.tx.Delete(taskEntry).Error
		if err != nil {
			return
		}
	}
	return
}

// DeleteTaskEntriesByTaskID deletes TaskEntry records by task's ID
func (repo *TaskEntryRepository) DeleteTaskEntriesByTaskID(taskID string) (err error) {
	return repo.tx.Delete(&model.Task{}, "task_id = ?", taskID).Error
}

// DeleteTaskEntriesByBoardID deletes TaskEntry records by board's ID
func (repo *TaskEntryRepository) DeleteTaskEntriesByBoardID(boardID string) (err error) {
	return repo.tx.Delete(&model.Task{}, "board_id = ?", boardID).Error
}
