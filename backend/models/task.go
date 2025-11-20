package models

import (
	"time"
)

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" gorm:"default:'medium'"` // low, medium, high
	Status      string    `json:"status" gorm:"default:'todo'"`     // todo, in-progress, done
	Tags        string    `json:"tags"`                             // JSON array stored as string
	DueDate     *time.Time `json:"due_date"`
	IsQuickTask bool      `json:"is_quick_task" gorm:"default:false"`
	Quadrant    string    `json:"quadrant"` // urgent-important, not-urgent-important, urgent-not-important, not-urgent-not-important
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `json:"user_id"`
}
