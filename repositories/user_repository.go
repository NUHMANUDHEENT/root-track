package repositories

import (
	"errors"
	"roottrack-backend/config"
	"roottrack-backend/models"

	"github.com/google/uuid"
)

type UserRepository struct{}

func (r *UserRepository) CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}
