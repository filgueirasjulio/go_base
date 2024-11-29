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

// @Summary Registrar usuário (Step 1)
// @Description Cria um novo usuário sem senha
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body requests.RegisterRequestStep1 true "Dados do usuário"
// @Success 201 {object} resources.UserResource "Usuário registrado"
// @Router /api/register/step1 [post]
func (ac *AuthController) RegisterStep1(c *fiber.Ctx) error {
	req := requests.NewRegisterRequest(ac.controller.DB)

	if err := c.BodyParser(&req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, "Falha ao processar requisição")
	}

	if err := validator.New().Struct(req.RegisterRequestStep1); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, "Falha ao iniciar o step 1")
	}

	user, err := ac.authService.RegisterUserStep1(req.RegisterRequestStep1)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "Erro ao registrar o usuário")
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
// @Param req body requests.RegisterRequestStep2 true "Dados do usuário"
// @Success 200 {object} resources.UserResource "Senha definida"
// @Router /api/register/step2 [post]
func (ac *AuthController) RegisterStep2(c *fiber.Ctx) error {
    req := requests.NewRegisterRequest(ac.controller.DB)

    if err := c.BodyParser(&req.RegisterRequestStep2); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusBadRequest, "Falha ao processar requisição")
    }

    if err := validator.New().Struct(req.RegisterRequestStep2); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, "Falha ao iniciar o step 2")
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

// envio de códgio para validação do cadastro de usuário
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
