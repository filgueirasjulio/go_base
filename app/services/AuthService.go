package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"crypto/md5"
    "encoding/hex"

	"tradeapi/app/mail"
	"tradeapi/app/models"
	"tradeapi/app/repositories"
	 requestsAuth "tradeapi/app/requests/auth"
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
	hashRepo *repositories.RegisterHashRepository
}

type ResponseStep2 struct {
    User  *models.User
    Token string
}

// NewAuthService cria uma nova instância do AuthService com o repositório injetado
func NewAuthService(db *gorm.DB) *AuthService {
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	codeRepo := repositories.NewValidationCodeRepository(db)
	hashRepo := repositories.NewRegisterHashRepository(db)
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
		codeRepo: codeRepo,
		hashRepo: hashRepo,
	}
}

//Login do usuário
func (s *AuthService) Login(params requestsAuth.LoginRequestParams) (string, int, error) {
    // Verifica se o usuário existe
    user, err := s.userRepo.FindByEmail(params.Email)
    if err != nil || user.ID == 0 {
        return "", 404, errors.New("login - erro ao buscar usuário")
    }

    // Verifica senha
    if !helpers.VerifyPassword(user.Password, params.Password) {
        return "", 401, errors.New("usuário ou senha inválidos")
    }

    // Gera token
    token, err := s.GenerateToken(user.Email, params.Password)
    if err != nil {
        return "", 500, errors.New("login - erro ao gerar token")
    }

    return token, 200, nil
}

// RegisterUser registra um novo usuário
func (s *AuthService) RegisterUserStep1(req requestsAuth.RegisterRequestStep1) (*models.User, error) {
    // Verifica se o e-mail já está registrado
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("e-mail já registrado")
    }

    // Cria o novo usuário sem senha
    userData := &models.User{
        Name:  req.Name,
        Email: req.Email,
    }

    // Salva o usuário no banco
    user, err := s.authRepo.Create(userData)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthService) RegisterUserStep2(req requestsAuth.RegisterRequestStep2) (*ResponseStep2, int, error) {
	registerHash, err := s.hashRepo.GetLastHash(req.Hash)

    if err != nil && err.Error() != "record not found" {
        return nil, 500, errors.New("erro ao buscar o código de validação no banco")
    }

    if registerHash == nil || registerHash.Hash == "" || registerHash.ExpiresAt.IsZero() {
        return nil, 404, errors.New("código de validação não encontrado")
    }

    if registerHash.Hash != req.Hash {
        return nil, 404, errors.New("código inválido")
    }

    if time.Now().After(registerHash.ExpiresAt) {
        return nil, 403, errors.New("código expirado, gere um novo")
    }

    user, err := s.userRepo.FindByHash(req.Hash)
    if err != nil {
        return nil, 500, errors.New("error ao buscar usuário")
    }
    if user == nil {
        return nil, 404, errors.New("usuário não encontrado")
    }

    // Cria o hash da senha
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, 500, errors.New("falha ao gerar o hash da senha")
    }

    // Atualiza a senha do usuário
    user.Password = string(hashedPassword)
    // Marca como ativo no sistema
    user.IsActive = true

    // Salva as alterações
    _, err = s.userRepo.Update(user)
    if err != nil {
        return nil, 500, errors.New("erro au atualizar o password do usuário")
    }

	token, err := s.GenerateToken(user.Email, user.Password)
    if err != nil {
        return nil, 500, errors.New("erro ao gerar token")
	}

    return &ResponseStep2{User: user, Token: token}, 200, nil
}

// envio de códgio para validação do cadastro de usuário
func (s *AuthService) SendValidationCodeEmail(email string) error {

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuário não encontrado")
	}

	code := helpers.GenerateValidationCode()
	validationCodeData := &models.ValidationCode{
		UserEmail:  email,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	validationCode, err := s.codeRepo.Create(validationCodeData)
	if err != nil {
		return errors.New("falha ao registrar código de validação")
	}

	mail.NewValidationCodeMail(user.Email, user.Name, validationCode.Code, user.IsActive,)

	return nil
}

// valida o código de ativação.
func (s *AuthService) VerifyValidationCode(c *fiber.Ctx, email, code string) (string, int, error) {
    validationCode, err := s.codeRepo.GetLastCodeByUserEmail(email)

    if err != nil && err.Error() != "record not found" {
        return "", 500, errors.New("erro ao buscar o código de validação no banco")
    }

    if validationCode == nil || validationCode.Code == "" || validationCode.ExpiresAt.IsZero() {
        return "", 404, errors.New("código de validação não encontrado")
    }

    if code != validationCode.Code {
        return "", 404, errors.New("código inválido")
    }

    if time.Now().After(validationCode.ExpiresAt) {
        return "", 403, errors.New("código expirado, gere um novo")
    }

    // Marca o código como validado e deleta
    err = s.codeRepo.DeleteCodeByUserEmail(email)
    if err != nil {
        return "", 500, fmt.Errorf("Falha ao deletar o código de validação")
    }

    // Gera hash MD5
    hash := md5.New()
    hash.Write([]byte(email + time.Now().String()))
    hashMD5 := hex.EncodeToString(hash.Sum(nil))

    // Cria registro de hash temporário
    registerHash := models.RegisterHash{
        UserEmail:      email,
        Hash:       hashMD5,
        CreatedAt:  time.Now(),
        ExpiresAt:  time.Now().Add(30 * time.Minute),
    }

    // Salva hash temporário
    _, err = s.authRepo.CreateHash(&registerHash)
    if err != nil {
        return "", 500, errors.New("Falha ao salvar hash temporário")
    }

    return hashMD5, 200, nil
}

func (s *AuthService) GenerateToken(email string, senha string) (string, error) {
    // Tempo de expiração
    expirationTime := time.Now().Add(72 * time.Hour) // 3 dias

    // Configuração das claims
    claims := jwt.MapClaims{
        "email": email,
        "exp":   expirationTime.Unix(),
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
        return "", errors.New("erro ao gerar token")
    }

    return tokenString, nil
}