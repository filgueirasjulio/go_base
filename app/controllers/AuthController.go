package controllers

import (
	"log"

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

// @Summary Registrar usuário
// @Description Cria um novo usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body requests.RegisterRequest true "Dados do usuário"
// @Success 201 {object} resources.UserResource "Usuário registrado"
// @Router /api/register [post]
func (ac *AuthController) Register(c *fiber.Ctx) error {
	req := requests.NewRegisterRequest(ac.controller.DB)
	if req.DB == nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "DB não inicializado")
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, "Falha ao processar requisição")
	}

	errMsg := req.Validate()
	if errMsg != nil {
		helpers.ErrorResponse(c, fiber.StatusUnprocessableEntity, errMsg)
		return nil
	}

	user, err := ac.authService.RegisterUser(*req)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "Erro ao registrar o usuário")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Usuário registrado com sucesso!",
		"data":    resources.Transform(user),
	})
}

// envio de códgio para validação do cadastro de usuário
func (ac *AuthController) SendValidationCode(c *fiber.Ctx) error {

	var dados map[string]string
	err := c.BodyParser(&dados)
	userID := dados["user_id"]

	err = ac.authService.SendValidationCodeEmail(userID)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SuccessResponseString(c, fiber.StatusOK, "Código enviado com sucesso")
}

func (ac *AuthController) VerifyValidationCode(c *fiber.Ctx) error {
    var dados map[string]string
    err := c.BodyParser(&dados)
    userID := dados["user_id"]
    code := dados["code"]

    token, status, err := ac.authService.VerifyValidationCode(c, userID, code)
    if err != nil {
        log.Println("Erro no serviço:", err)
        return c.Status(status).JSON(fiber.Map{"error": err.Error()})
    }
    // Log do token antes de retornar
    log.Println("Token gerado:", token)

    return c.Status(status).JSON(fiber.Map{
        "token": token,
    })
}
