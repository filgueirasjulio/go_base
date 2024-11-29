package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"tradeapi/app/requests"
	"tradeapi/app/resources"
	"tradeapi/app/services"
	"tradeapi/utils/helpers"
)

type AuthController struct {
	controller  *Controller
	authService *services.AuthService
}

func NewAuthController(controller *Controller) *AuthController {
	return &AuthController{
		controller:  controller,
		authService: services.NewAuthService(controller.DB),
	}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	req := requests.NewLoginRequest(ac.controller.DB)

	if err := c.BodyParser(&req.LoginRequestParams); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	token, status, err := ac.authService.Login(req.LoginRequestParams)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(status).JSON(fiber.Map{
		"token": token,
	})
}

// @Summary Registrar usuário (Step 1)
// @Description Cria um novo usuário sem senha
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body requests.RegisterRequestStep1 true "Dados do usuário"
// @Success 201 {object} resources.UserResource "Usuário registrado"
// @Router /api/register/step1 [post]
func (ac *AuthController) RegisterStep1(c *fiber.Ctx) error {
	req := requests.NewRegisterRequest(ac.controller.DB)

	if err := c.BodyParser(&req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error())
	}

	if err := validator.New().Struct(req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	user, err := ac.authService.RegisterUserStep1(req.RegisterRequestStep1)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error())
	}

	//envio de e-mail
	err = ac.authService.SendValidationCodeEmail(user.Email)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error())
	}


	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Foi enviado um código de validação para seu e-mail. Você tem até 30 minutos para o validar",
		"data":    resources.Transform(user),
	})
}

// @Summary Registrar usuário (Step 2)
// @Description Define a senha do usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body requests.RegisterRequestStep2 true "Dados do usuário"
// @Success 200 {object} resources.UserResource "Senha definida"
// @Router /api/register/step2 [post]
func (ac *AuthController) RegisterStep2(c *fiber.Ctx) error {
    req := requests.NewRegisterRequest(ac.controller.DB)

    if err := c.BodyParser(&req.RegisterRequestStep2); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error())
    }

    if err := validator.New().Struct(req.RegisterRequestStep2); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error())
    }

    response, status, err := ac.authService.RegisterUserStep2(req.RegisterRequestStep2)
    if err != nil {
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Senha definida com sucesso!",
		"token": response.Token,
        "data":    resources.Transform(response.User),
    })
}

// @Summary Envia código de validação
// @Description Envia código de validação para o e-mail fornecido.
// @Tags Auth
// @Param email body string true "E-mail do usuário" example="joão@mail.com"
// @Success 200 {object} map[string]string "Código enviado" example="{\"message\": \"Código enviado com sucesso\"}"
// @Router /api/send-validation-code [post]
func (ac *AuthController) SendValidationCode(c *fiber.Ctx) error {

	var dados map[string]string
	err := c.BodyParser(&dados)
	email := dados["email"]

	err = ac.authService.SendValidationCodeEmail(email)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SuccessResponseString(c, fiber.StatusOK, "Código enviado com sucesso")
}

// @Summary Verifica código de validação
// @Description Verifica se o código de validação é válido.
// @Tags Auth
// @Param email body string true "E-mail do usuário" example="joao@mail.com"
// @Param code body string true "Código de validação" example="5329"
// @Success 200 {object} map[string]string "Código válido"
// @Router /api/verify-validation-code [post]
func (ac *AuthController) VerifyValidationCode(c *fiber.Ctx) error {
	var dados map[string]string
	err := c.BodyParser(&dados)
	email := dados["email"]
	code := dados["code"]

	hash, status, err := ac.authService.VerifyValidationCode(c, email, code)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(status).JSON(fiber.Map{
		"hash": hash,
	})
}
