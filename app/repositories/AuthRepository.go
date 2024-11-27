package repositories

import (
	"tradeapi/app/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository cria uma nova instância do AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// Create salva um novo usuário no banco de dados
func (r *AuthRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}