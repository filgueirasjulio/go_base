package services

import (
	"errors"
	"time"
	"gorm.io/gorm"

	"tradeapi/app/models"
	"tradeapi/app/requests"
	"tradeapi/app/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

// NewAuthService cria uma nova instância do AuthService com o repositório injetado
func NewAuthService(db *gorm.DB) *AuthService {
	authRepo := repositories.NewAuthRepository(db)
	return &AuthService{
		authRepo: authRepo,
	}
}

// RegisterUser registra um novo usuário
func (s *AuthService) RegisterUser(req requests.RegisterRequestInternal) (*models.User, error) {
	// Cria o hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("falha ao gerar o hash da senha")
	}

	// Cria o novo usuário
	userData := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Salva o usuário no banco
	user, err := s.authRepo.Create(userData)
	if err != nil {
		return nil, err
	}

	return user, nil
}
