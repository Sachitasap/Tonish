package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" gorm:"default:'medium'"` // low, medium, high
	Status      string     `json:"status" gorm:"default:'todo'"`     // todo, in-progress, done
	Tags        string     `json:"tags"`                             // JSON array stored as string
	DueDate     *time.Time `json:"due_date"`
	IsQuickTask bool       `json:"is_quick_task" gorm:"default:false"`
	Quadrant    string     `json:"quadrant"` // urgent-important, not-urgent-important, urgent-not-important, not-urgent-not-important
	TaskType    string     `json:"task_type" gorm:"default:'kanban'"`
	IsArchived  bool       `json:"is_archived" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uint       `json:"user_id"`
}
