package controllers

import (
	"github.com/gofiber/fiber/v2"
	"tradeapi/app/resources"
	"tradeapi/app/services"
	"tradeapi/utils/helpers"
)

type UserController struct {
	controller  *Controller
	userService *services.UserService
}

// Construtor do UserController com injeção do UserService
func NewUserController(controller *Controller) *UserController {
	return &UserController{
		controller:  controller,
		userService: services.NewUserService(controller.DB), 
	}
}

// @Summary Listar usuários
// @Description Obtém uma lista de usuários cadastrados
// @Tags Users
// @Success 200 {array} models.User
// @Router /api/users [get]
func (uc *UserController) Index(c *fiber.Ctx) error {
	// Obtém os usuários usando o serviço
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "Erro ao listar usuários")
	}

	// Transforma os usuários em recursos
	userResources := resources.TransformCollection(users)

	// Retorna os dados formatados como JSON
	return c.JSON(userResources)
}
