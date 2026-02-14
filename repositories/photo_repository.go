package repositories

import (
	"roottrack-backend/config"
	"roottrack-backend/models"

	"github.com/google/uuid"
)

type PhotoRepository struct{}

func (r *PhotoRepository) Create(photo *models.ProgressPhoto) error {
	return config.DB.Create(photo).Error
}

func (r *PhotoRepository) GetAllByUser(userID uuid.UUID) ([]models.ProgressPhoto, error) {
	var photos []models.ProgressPhoto
	err := config.DB.Where("user_id = ?", userID).Order("taken_at desc").Find(&photos).Error
	return photos, err
}

func (r *PhotoRepository) Delete(id uuid.UUID) error {
	return config.DB.Delete(&models.ProgressPhoto{}, id).Error
}

func (r *PhotoRepository) GetLatestByUser(userID uuid.UUID) (*models.ProgressPhoto, error) {
	var photo models.ProgressPhoto
	err := config.DB.Where("user_id = ?", userID).Order("taken_at desc").First(&photo).Error
	return &photo, err
}
