package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Routine struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	ProductID  *uuid.UUID     `gorm:"type:uuid" json:"product_id"`
	Title      string         `json:"title"`
	StartDate  time.Time      `json:"start_date"`
	EndDate    *time.Time     `json:"end_date"`
	Recurrence string         `json:"recurrence"`                      // none, daily, weekly, custom
	DaysOfWeek pq.StringArray `gorm:"type:text[]" json:"days_of_week"` // postgres array
	Completed  bool           `json:"completed"`
	Notes      string         `json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`

	Product *Product `gorm:"foreignKey:ProductID;references:ID"`
}

func (r *Routine) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
