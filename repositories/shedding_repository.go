package repositories

import (
	"roottrack-backend/config"
	"roottrack-backend/models"

	"github.com/google/uuid"
)

type SheddingRepository struct{}

func (r *SheddingRepository) Create(log *models.SheddingLog) error {
	return config.DB.Create(log).Error
}

func (r *SheddingRepository) GetAllByUser(userID uuid.UUID) ([]models.SheddingLog, error) {
	var logs []models.SheddingLog
	err := config.DB.Where("user_id = ?", userID).Order("date desc").Find(&logs).Error
	return logs, err
}

func (r *SheddingRepository) GetWeeklyAverage(userID uuid.UUID) (float64, error) {
	var avg float64
	err := config.DB.Model(&models.SheddingLog{}).
		Where("user_id = ? AND date >= now() - interval '7 days'", userID).
		Select("AVG(shedding_count)").
		Scan(&avg).Error
	return avg, err
}
