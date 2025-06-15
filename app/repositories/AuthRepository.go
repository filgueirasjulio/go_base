package repositories

import (
	"base/app/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db       *gorm.DB
	userRepo *UserRepository
}

// NewAuthRepository cria uma nova instância do AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db,
		userRepo: NewUserRepository(db),
	}
}

// Create salva um novo usuário no banco de dados
func (r *AuthRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// cria o hash após validação do token de 4 digitos
func (r *AuthRepository) CreateHash(registerHash *models.RegisterHash) (string, error) {

	err := r.db.Create(&registerHash).Error
	if err != nil {
		return "", err
	}

	return "", nil
}
