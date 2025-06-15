package services

import (
	"strconv"
	"base/app/models"
	"base/app/repositories"
	requestsUser "base/app/requests/user"
	"base/utils/helpers"

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

// Método para obter um usuário em específico
func (s *UserService) FindById(id string) (*models.User, error) {
	userId, _ := strconv.Atoi(id)

	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Método para atualizar detalhes de um usuário em específico
func (s *UserService) UpdateDetails(id int, req requestsUser.UserUpdateDetailsRequestParams) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	birthDate, err := helpers.ConvertData(req.BirthDate)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Phone = helpers.RemovePhoneMask(req.Phone) 
	user.BirthDate = &birthDate

	user, err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
