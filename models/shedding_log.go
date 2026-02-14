package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SheddingLog struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Date          time.Time `json:"date"`
	SheddingCount int       `json:"shedding_count"`
	CreatedAt     time.Time `json:"created_at"`
}

func (s *SheddingLog) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
