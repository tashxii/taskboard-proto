package model

import (
	"taskboard/common"
	"time"
)

// Task present task of the app.
type Task struct {
	ID             string    `gorm:"primary_key;size:32"`
	Name           string    `gorm:"not null;size:255"`
	Description    string    `gorm:"size:8000"`
	AssigneeUserID string    `gorm:"size:32"`
	BoardID        string    `gorm:"size:32"`
	CreatedDate    time.Time `gorm:"not null"`
	IsClosed       bool      `gorm:"not null"`
	Version        int       `gorm:"not null"` // Version for optimistic lock
	EsitmateSize   int
}

// NewTask returns created new task
func NewTask(name, description string, isClosed bool, now time.Time) *Task {
	return &Task{
		ID:          "task_" + common.GenerateID(),
		Name:        name,
		Description: description,
		IsClosed:    isClosed,	
		CreatedDate: now,
		Version:     1,
	}
}
