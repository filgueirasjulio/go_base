package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"tradeapi/app/mail"
	"tradeapi/app/models"
	"tradeapi/app/repositories"
	"tradeapi/app/requests"
	"tradeapi/utils/helpers"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Token string `json:"token"`
}

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
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Salva o usuário no banco
	user, err := s.authRepo.Create(userData)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// envio de códgio para validação do cadastro de usuário
func (s *AuthService) SendValidationCodeEmail(id string) error {

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
		UserID:    uint(userID),
		Code:      code,
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

// valida o código de ativação.
func (s *AuthService) VerifyValidationCode(c *fiber.Ctx, user_id, code string) (string, int, error) {
	validationCode, err := s.codeRepo.GetLastCodeByUserId(user_id)

	if err != nil && err.Error() != "record not found" {
		return "", 500, fmt.Errorf("erro ao buscar o código de validação no banco")
	}

	if validationCode == nil || validationCode.Code == "" || validationCode.ExpiresAt.IsZero() {
		return "", 404, fmt.Errorf("código de validação não encontrado")
	}

	if code != validationCode.Code {
		return "", 404, fmt.Errorf("código inválido")
	}

	if time.Now().After(validationCode.ExpiresAt) {
		return "", 403, fmt.Errorf("código expirado, gere um novo")
	}

	// Marca o código como validado e salva
	err = s.codeRepo.MarkAsValidated(validationCode)
	if err != nil {
		return "", 500, fmt.Errorf("Falha ao atualizar o status do código de validação")
	}

	// Se chegamos até aqui, podemos gerar o token
	token, err := s.GenerateToken(user_id)
	if err != nil {
		return "", 500, fmt.Errorf("Falha ao gerar token")
	}

	// Retorna o token gerado
	return token, 200, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	// Tempo de expiração
	expirationTime := time.Now().Add(72 * time.Hour) // 3 dias

	// Configuração das claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	// Criação do token com método de assinatura HS256 (HMAC com SHA-256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Chave secreta
	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "JWT_KEY não configurada")
	}

	// Geração do token com assinatura
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
