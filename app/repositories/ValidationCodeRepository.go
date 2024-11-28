package repositories

import (
	"tradeapi/app/models"
	"fmt"
	"gorm.io/gorm"
)

type ValidationCodeRepository struct {
	db *gorm.DB
}

// NewValidationCodeRepository cria uma nova instância do ValidationCodeRepository
func NewValidationCodeRepository(db *gorm.DB) *ValidationCodeRepository {
	return &ValidationCodeRepository{db: db}
}

// Obtem o último código de validação gerado pelo usuário
func (r *ValidationCodeRepository) GetLastCodeByUserId(id string) (*models.ValidationCode, error) {
    var code models.ValidationCode
    result := r.db.Where("user_id = ?", id).Order("created_at DESC").First(&code)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error  
    }
    return &code, nil
}

// Create salva um novo usuário no banco de dados
func (r *ValidationCodeRepository) Create(validationCode *models.ValidationCode) (*models.ValidationCode, error) {
    if err := r.db.Create(validationCode).Error; err != nil {
        return nil, fmt.Errorf("erro ao salvar o código de validação: %w", err)
    }
    return validationCode, nil
}

//atualizado o codigo como validado
func (r *ValidationCodeRepository) MarkAsValidated(validationCode *models.ValidationCode) error {
    validationCode.IsValidated = true

    result := r.db.Save(validationCode)
    if result.Error != nil {
        return fmt.Errorf("erro ao marcar o código como validado: %w", result.Error)
    }
    return nil
}