package repositories

import (
	"base/app/models"
	"gorm.io/gorm"
)

type RegisterHashRepository struct {
	db *gorm.DB
}

// NewRegisterHashRepository cria uma nova instância do RegisterHashRepository
func NewRegisterHashRepository(db *gorm.DB) *RegisterHashRepository {
	return &RegisterHashRepository{db: db}
}


func (r *RegisterHashRepository) FindByHash(hashStr string) (*models.RegisterHash, error) {
    var hash models.RegisterHash
    result := r.db.Where("hash = ?", hashStr).Order("created_at DESC").First(&hash)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error  
    }
    return &hash, nil
}

// Obtem o último código de validação gerado pelo usuário
func (r *RegisterHashRepository) GetLastHash(email string) (*models.RegisterHash, error) {
    var hash models.RegisterHash
    result := r.db.Where("hash = ?", email).Order("created_at DESC").First(&hash)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error  
    }
    return &hash, nil
}

//obtem registro via rash
func (r *RegisterHashRepository) FindUserByHash(hash string) (*models.RegisterHash, error) {
    register, err := r.FindByHash(hash)
    if err != nil {
        return nil, err
    }
    return register, nil
}