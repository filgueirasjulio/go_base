package repositories

import (
	"tradeapi/app/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	hashRepo *RegisterHashRepository
}

// Construtor para UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	hashRepo := NewRegisterHashRepository(db)
	return &UserRepository{db: db,
		hashRepo: hashRepo,
	}
}

//Método para obter usuário via id
func (r *UserRepository) FindByID(id int) (*models.User, error) {
    user := &models.User{}
    return user, r.db.First(user, id).Error
}

//busca o registro via e-mail
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//busca o registro via hash gerado no processo de cadastro
func (r *UserRepository) FindByHash(hash string) (*models.User, error) {
	registerHash, err := r.hashRepo.FindByHash(hash)
	if err != nil {
		return nil, err
	}

	return r.FindByEmail(registerHash.UserEmail)
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

// Update atualiza um usuário existente
func (r *UserRepository) Update(user *models.User) (*models.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
