package repositories

import (
	"roottrack-backend/config"
	"roottrack-backend/models"

	"github.com/google/uuid"
)

type ProductRepository struct{}

func (r *ProductRepository) Create(product *models.Product) error {
	return config.DB.Create(product).Error
}

func (r *ProductRepository) GetAllByUser(userID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Where("user_id = ?", userID).Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *models.Product) error {
	return config.DB.Save(product).Error
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	return config.DB.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := config.DB.First(&product, id).Error
	return &product, err
}
