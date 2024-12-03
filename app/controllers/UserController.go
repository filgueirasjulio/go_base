package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"tradeapi/app/resources"
	"tradeapi/app/services"
	"tradeapi/utils/helpers"
	"strconv"
	requestsUser "tradeapi/app/requests/user"
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

// @Summary Listar usuário
// @Description Obtém informações de um usuário em específico
// @Tags Users
// @Param id path string true "ID do usuário"
// @Success 200 {array} models.User
// @Router /api/users/{id} [get]
func (uc *UserController) Show(c *fiber.Ctx) error {
    id := c.Params("id")

    user, err := uc.userService.FindById(id)
    if err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "Erro ao listar usuário")
    }

    userResources := resources.Transform(user)

    return c.JSON(userResources)
}

// @Summary Atualizar usuário
// @Description Atualiza informações de um usuário em específico
// @Tags Users
// @Param request body requestsAuth.LoginRequestParams true "Dados do usuário"
// @Success 200 {array} models.User
// @Router /api/users/{id} [put]
func (uc *UserController) UpdateDetails(c *fiber.Ctx) error {
	req := requestsUser.NewUserUpdateDetailsRequest(uc.controller.DB)

    if err := c.BodyParser(&req.UserUpdateDetailsRequestParams); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error())
    }

    if err := validator.New().Struct(req.UserUpdateDetailsRequestParams); err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error())
    }

	userId, _ := strconv.Atoi(c.Params("id"))

    user, err := uc.userService.UpdateDetails(userId, req.UserUpdateDetailsRequestParams)
    if err != nil {
        return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, "Erro ao atualizar usuário")
    }

    userResources := resources.Transform(user)

    return c.JSON(userResources)
}

