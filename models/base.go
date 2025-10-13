package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp(6);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp(6);autoCreateTime;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp(6)" json:"deleted_at"`

	// CreatedBy uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by"`
	// UpdatedBy uuid.UUID `gorm:"type:uuid;column:updated_by;null" json:"updated_by"`
	// DeletedBy uuid.UUID `gorm:"type:uuid;column:deleted_by;null" json:"deleted_by"`

	// UserCreatedBy *User `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	// UserUpdatedBy *User `gorm:"foreignKey:UpdatedBy;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	// UserDeletedBy *User `gorm:"foreignKey:DeletedBy;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

type PaginationRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
