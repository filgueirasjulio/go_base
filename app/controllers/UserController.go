package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
	requestsUser "base/app/requests/user"
	"base/app/resources"
	"base/app/services"
	"base/utils/helpers"
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
	start := time.Now()

	users, err := uc.userService.GetAllUsers()
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, helpers.ErrUsersIndex, uc.controller.Logger)
	}

	if users == nil {
		return helpers.SuccessResponseData(c, 200, map[string]interface{}{"users": []interface{}{}}, helpers.ErrUsersNotFound, uc.controller.Logger, start)
	}

	// Transforma os usuários em recursos
	userResources := resources.TransformCollection(users)

	// Converte para JSON com a ordem correta dos campos
	data := map[string]interface{}{
		"users": userResources,
	}

	// Retorna a resposta
	return helpers.SuccessResponseData(c, 200, data, helpers.UserIndex, uc.controller.Logger, start)
}

// @Summary Listar usuário
// @Description Obtém informações de um usuário em específico
// @Tags Users
// @Param id path string true "ID do usuário"
// @Success 200 {array} models.User
// @Router /api/users/{id} [get]
func (uc *UserController) Show(c *fiber.Ctx) error {
	start := time.Now()
	id := c.Params("id")

	user, err := uc.userService.FindById(id)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, helpers.ErrUserShow, uc.controller.Logger)
	}

	// Transformação do usuário
	transformedUser := resources.Transform(user)

	// Construção do mapa de resposta
	data := map[string]interface{}{
		"users": transformedUser,
	}

	return helpers.SuccessResponseData(c, 200, data, helpers.UserShow, uc.controller.Logger, start)
}

// @Summary Atualizar usuário
// @Description Atualiza informações de um usuário em específico
// @Tags Users
// @Param request body requestsUser.UserUpdateDetailsRequestParams true "Dados do usuário"
// @Success 200 {array} models.User
// @Router /api/users/{id} [put]
func (uc *UserController) UpdateDetails(c *fiber.Ctx) error {
	start := time.Now()

	req := requestsUser.NewUserUpdateDetailsRequest(uc.controller.DB)

	if err := c.BodyParser(&req.UserUpdateDetailsRequestParams); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusBadRequest, err.Error(), uc.controller.Logger)
	}

	if err := validator.New().Struct(req.UserUpdateDetailsRequestParams); err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusUnprocessableEntity, err.Error(), uc.controller.Logger)
	}

	userId, _ := strconv.Atoi(c.Params("id"))

	user, err := uc.userService.UpdateDetails(userId, req.UserUpdateDetailsRequestParams)
	if err != nil {
		return helpers.ErrorResponseString(c, fiber.StatusInternalServerError, helpers.ErrUserUpdate, uc.controller.Logger)
	}

	userResources := resources.Transform(user)

	data := map[string]interface{}{
		"users": c.JSON(userResources),
	}

	return helpers.SuccessResponseData(c, 200, data, helpers.UserUpdate, uc.controller.Logger, start)
}
