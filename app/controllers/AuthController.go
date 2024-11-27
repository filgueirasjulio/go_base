package controllers

import (
	"github.com/gofiber/fiber/v2"

	"tradeapi/utils/helpers"
	"tradeapi/app/requests"
	"tradeapi/app/resources"
	"tradeapi/app/services"
)

type AuthController struct {
	controller *Controller
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