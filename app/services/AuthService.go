package services

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	"tradeapi/app/models"
	"tradeapi/app/repositories"
	"tradeapi/app/requests"
	"tradeapi/utils/helpers"
	"tradeapi/app/mail"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
	userRepo *repositories.UserRepository
	codeRepo *repositories.ValidationCodeRepository
}

// NewAuthService cria uma nova instância do AuthService com o repositório injetado
func NewAuthService(db *gorm.DB) *AuthService {
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	codeRepo := repositories.NewValidationCodeRepository(db)
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
		codeRepo: codeRepo,
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

//envio de códgio para validação do cadastro de usuário
func (s *AuthService) SendVerificationCodeEmail(id string) error {

	userID, err := strconv.Atoi(id)
    if err != nil {
        return errors.New("ID inválido")
    }

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("usuário não encontrado")
    }

	code := helpers.GenerateValidationCode()
    validationCodeData := &models.ValidationCode{
        UserID:       uint(userID),
        Code:         code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
    }

	validationCode, err := s.codeRepo.Create(validationCodeData)
	if err != nil {
		return errors.New("falha ao registrar código de validação")
	}

	mail.NewValidationCodeMail(user.Email, user.Name, validationCode.Code)

    return nil
}