package repositories

import (
	"tradeapi/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// Construtor para UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Método para obter todos os usuários
func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
