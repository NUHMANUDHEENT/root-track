package repositories

import (
	"roottrack-backend/config"
	"roottrack-backend/models"
	"time"

	"github.com/google/uuid"
)

type RoutineRepository struct{}

func (r *RoutineRepository) Create(routine *models.Routine) error {
	return config.DB.Create(routine).Error
}

func (r *RoutineRepository) GetAllByUser(userID uuid.UUID) ([]models.Routine, error) {
	var routines []models.Routine
	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&routines).Error

	return routines, err
}
func (r *RoutineRepository) GetActiveByUser(userID uuid.UUID, today time.Time) ([]models.Routine, error) {
	var routines []models.Routine

	err := config.DB.
		Preload("Product").
		Joins("LEFT JOIN products ON products.id = routines.product_id").
		Where("routines.user_id = ?", userID).
		Where("routines.start_date <= ?", today).
		Where("routines.end_date IS NULL OR routines.end_date >= ?", today).
		Where("(routines.product_id IS NULL OR products.is_active = ?)", true).
		Find(&routines).Error

	return routines, err
}

func (r *RoutineRepository) FindByID(id uuid.UUID) (*models.Routine, error) {
	var routine models.Routine
	err := config.DB.First(&routine, id).Error
	return &routine, err
}

func (r *RoutineRepository) Update(routine *models.Routine) error {
	return config.DB.Save(routine).Error
}

func (r *RoutineRepository) Delete(id uuid.UUID) error {
	return config.DB.Delete(&models.Routine{}, id).Error
}
