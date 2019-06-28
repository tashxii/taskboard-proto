package model

import (
	"math"
	"taskboard/common"
	"time"
)

// Board presents a board which has plural tasks
type Board struct {
	ID          string       `gorm:"primary_key;size:32"`
	Name        string       `gorm:"unique;size:255"`
	DispOrder   int          `gorm:"not null"`
	IsSystem   bool         `gorm:"not null"`
	IsClosed    bool         `gorm:"not null"`
	CreatedDate time.Time    `gorm:"not null"`
	Version     int          `gorm:"not null"` // Version for optimistic lock
	TaskEntries []*TaskEntry `gorm:"-"`        // Ignore, fetch by service layer when they needed
}

// NewBoard returns created new board
func NewBoard(name string, isSystem, isClosed bool, now time.Time) *Board {
	return &Board{
		ID:          "board_" + common.GenerateID(),
		Name:        name,
		DispOrder:   math.MaxInt32,
		IsSystem:   isSystem,
		IsClosed:    isClosed,
		CreatedDate: now,
		Version:     1,
		TaskEntries: []*TaskEntry{},
	}
}

// AddTaskEntry adds & returns created task entry
func (board *Board) AddTaskEntry(taskID string) *TaskEntry {
	taskEntry := &TaskEntry{
		ID:        "taskentry_" + common.GenerateID(),
		DispOrder: len(board.TaskEntries),
		BoardID:   board.ID,
		TaskID:    taskID,
	}
	board.TaskEntries = append(board.TaskEntries, taskEntry)
	return taskEntry
}

// RemoveAllTaskEntries clears task entries of board
func (board *Board) RemoveAllTaskEntries() {
	board.TaskEntries = []*TaskEntry{}
}

// SetTaskEntries changes order of tasks
func (board *Board) SetTaskEntries(taskEntries []*TaskEntry) {
	for i, taskEntry := range taskEntries {
		taskEntry.DispOrder = i
	}
	board.TaskEntries = taskEntries
}
