package model

// TaskEntry presents entry of board
type TaskEntry struct {
	ID        string `gorm:"primary_key;size:32"`
	DispOrder int    `gorm:"not null"`
	BoardID   string `gorm:"not null;size:32"`
	TaskID    string `gorm:"not null;size:32"`
}
