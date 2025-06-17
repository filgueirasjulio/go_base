package controllers

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	requestsAuth "base/app/requests/auth"
	"base/app/resources"
	"base/app/services"
	"base/utils/helpers"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type AuthController struct {
	controller  *Controller
	authService *services.AuthService
}

func NewAuthController(controller *Controller) *AuthController {
	return &AuthController{
		controller:  controller,
		authService: services.NewAuthService(controller.DB, controller.Logger),
	}
}

// @Summary Logar usuário
// @Description Realiza login de usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param request body requestsAuth.LoginRequestParams true "Dados do usuário"
// @Success 201 {object} TokenResponse "Resposta de login"
// @Router /api/auth/login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
	start := time.Now()
	req := requestsAuth.NewLoginRequest(ac.controller.DB)

	if err := c.BodyParser(&req.LoginRequestParams); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error(), ac.controller.Logger)
	}

	if err := req.Validate(); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnprocessableEntity, err, ac.controller.Logger)
	}

	token, status, err := ac.authService.Login(req.LoginRequestParams)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	ac.controller.Logger.Info().Str("status", "sucesso").Int("codigo", 200).Msg(helpers.LoginSuccessfull)

	data := map[string]interface{}{
		"token": token,
	}

	return helpers.SuccessResponseData(c, status, data, helpers.LoginSuccessfull, ac.controller.Logger, start)
}

// @Summary Registrar usuário (Step 1)
// @Description Cria um novo usuário sem senha
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param request body requestsAuth.RegisterRequestStep1 true "Dados do usuário"
// @Success 201 {object} resources.UserResource "Usuário registrado"
// @Router /api/auth/register_step1 [post]
func (ac *AuthController) RegisterStep1(c *fiber.Ctx) error {
	start := time.Now()

	req := requestsAuth.NewRegisterRequest(ac.controller.DB)

	if err := c.BodyParser(&req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error(), ac.controller.Logger)
	}

	if err := validator.New().Struct(req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error(), ac.controller.Logger)
	}

	user, err := ac.authService.RegisterUserStep1(req.RegisterRequestStep1)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	//envio de e-mail
	err = ac.authService.SendValidationCodeEmail(user.Email)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	data := map[string]interface{}{
		"message": helpers.ValidationCodeSent,
		"user":    resources.Transform(user),
	}

	return helpers.SuccessResponseData(c, 200, data, helpers.RegisterFirstStep, ac.controller.Logger, start)
}

// @Summary Registrar usuário (Step 2)
// @Description Define a senha do usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param request body requestsAuth.RegisterRequestStep2 true "Dados do usuário"
// @Success 200 {object} resources.UserResource "Senha definida"
// @Router /api/auth/register_step2 [post]
func (ac *AuthController) RegisterStep2(c *fiber.Ctx) error {
	start := time.Now()
	req := requestsAuth.NewRegisterRequest(ac.controller.DB)

	if err := c.BodyParser(&req.RegisterRequestStep2); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error(), ac.controller.Logger)
	}

	if err := validator.New().Struct(req.RegisterRequestStep2); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error(), ac.controller.Logger)
	}

	response, status, err := ac.authService.RegisterUserStep2(req.RegisterRequestStep2)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error(), ac.controller.Logger)
	}

	data := map[string]interface{}{
		"message": helpers.PasswordDefined,
		"token":   response.Token,
		"users":   resources.Transform(response.User),
	}

	return helpers.SuccessResponseData(c, status, data, helpers.RegisterSecondStep, ac.controller.Logger, start)
}

// @Summary Envia código de validação
// @Description Envia código de validação para o e-mail fornecido.
// @Tags Autenticação
// @Param request body requestsAuth.ValidationCodeRequestParams true "e-mail"
// @Success 200 {object} map[string]string "Código enviado" example="{\"message\": \"Código enviado com sucesso\"}"
// @Router /api/auth/send-validation-code [post]
func (ac *AuthController) SendValidationCode(c *fiber.Ctx) error {
	start := time.Now()

	var dados map[string]string
	err := c.BodyParser(&dados)
	email := dados["email"]

	err = ac.authService.SendValidationCodeEmail(email)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	return helpers.SuccessResponseString(c, fiber.StatusOK, helpers.ValidationCodeSentShort, ac.controller.Logger, start)
}

// @Summary Verifica código de validação
// @Description Verifica se o código de validação é válido.
// @Tags Autenticação
// @Param email body string true "E-mail do usuário" example="joao@mail.com"
// @Param request body requestsAuth.VerifyCodeRequestParams true "email, code"
// @Success 200 {object} object "Código validado com sucesso"
// @Router /api/auth/verify-validation-code [post]
func (ac *AuthController) VerifyValidationCode(c *fiber.Ctx) error {
	start := time.Now()
	var dados map[string]string
	err := c.BodyParser(&dados)
	email := dados["email"]
	code := dados["code"]

	hash, status, err := ac.authService.VerifyValidationCode(c, email, code)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	data := map[string]interface{}{
		"hash": hash,
	}
	return helpers.SuccessResponseData(c, status, data, helpers.ValidationCodeConfirmed, ac.controller.Logger, start)
}

// @Summary Deslogar usuário
// @Description    Desloga o usuário atual
// @Tags           Autenticação
// @Accept         json
// @Produce        json
// @Success 200 {object} object "Usuário deslogado com sucesso"
// @Router /api/auth/logout [get]
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	start := time.Now()
	c.Locals("user", nil)
	c.ClearCookie("token")
	return helpers.SuccessResponseString(c, fiber.StatusOK, helpers.LogoutSuccessfull, ac.controller.Logger, start)
}

func (ac *AuthController) RefreshToken(refreshTokenString string, c *fiber.Ctx) (string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return "", helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return "", helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	// Gere novo token
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	newAccessTokenString, err := newAccessToken.SignedString([]byte(os.Getenv("JWT_KEY")))

	return newAccessTokenString, err
}

// @Summary     Refresca o token de acesso
// @Description Retorna um novo token de acesso válido por 1 hora
// @Tags        Autenticação
// @Accept      application/json
// @Produce     application/json
// @Param       Authorization header string true "Token de refresh"
// @Success     200 {string} string "Novo token de acesso"
// @Router      /api/auth/refresh-token [post]
func (ac *AuthController) HandleRefreshToken(c *fiber.Ctx) error {
	start := time.Now()
	refreshTokenString := c.Get("Authorization")

	// Remova o prefixo "Bearer " se necessário
	if strings.HasPrefix(refreshTokenString, "Bearer ") {
		refreshTokenString = strings.TrimSpace(refreshTokenString[7:])
	}

	newAccessToken, err := ac.RefreshToken(refreshTokenString, c)

	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error(), ac.controller.Logger)
	}

	data := map[string]interface{}{
		"token": newAccessToken,
	}
	return helpers.SuccessResponseData(c, 200, data, helpers.TokenGenerated, ac.controller.Logger, start)
}
