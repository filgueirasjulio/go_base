package repositories

import (
	"tradeapi/app/models"
	"gorm.io/gorm"
)

type ValidationCodeRepository struct {
	db *gorm.DB
}

// NewValidationCodeRepository cria uma nova instância do ValidationCodeRepository
func NewValidationCodeRepository(db *gorm.DB) *ValidationCodeRepository {
	return &ValidationCodeRepository{db: db}
}

// Create salva um novo usuário no banco de dados
func (r *ValidationCodeRepository) Create(validationCode *models.ValidationCode) (*models.ValidationCode, error) {
	if err := r.db.Create(validationCode).Error; err != nil {
		return nil, err
	}
	return validationCode, nil
}