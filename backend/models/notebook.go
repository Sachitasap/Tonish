package models

import (
	"time"
)

type Notebook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Tags      string    `json:"tags"` // JSON array stored as string
	IsPinned  bool      `json:"is_pinned" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	Pages     []Page    `json:"pages" gorm:"foreignKey:NotebookID"`
}

type Page struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	NotebookID uint      `json:"notebook_id"`
	Title      string    `json:"title" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text"` // Rich-text JSON content from TipTap
	Tags       string    `json:"tags"`                     // JSON array stored as string
	IsPinned   bool      `json:"is_pinned" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
