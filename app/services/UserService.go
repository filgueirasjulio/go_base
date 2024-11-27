package services

import (
	"tradeapi/app/models"
	"tradeapi/app/repositories"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

// Construtor do UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(db),
	}
}

// Método para obter todos os usuários
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}
