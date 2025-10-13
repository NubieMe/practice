package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Todo struct {
	Title  string         `gorm:"column:title;size:255" json:"title" validate:"required"`
	Todo   pq.StringArray `gorm:"column:todo;type:text[]" json:"todo" validate:"required"`
	Check  pq.StringArray `gorm:"column:check;type:text[]" json:"check" validate:"required"`
	Images pq.StringArray `gorm:"column:images;type:text[]" json:"images"`

	UserID uuid.UUID `gorm:"type:uuid;column:user_id" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user"`

	Base
}

type TodoRequest struct {
	Title  string   `json:"title" validate:"required"`
	Todo   []string `json:"todo" validate:"required"`
	Check  []string `json:"check"`
	Images []string `json:"images"`

	UserID uuid.UUID `json:"user_id"`
}
