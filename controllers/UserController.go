// controllers/userController.go
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"tradeapi/services"
	"tradeapi/resources"
)

type UserController struct{
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// @Summary Listar usuários
// @Description Obtém uma lista de usuários cadastrados
// @Tags Users
// @Success 200 {array} models.User
// @Router /api/users [get]
func (uc *UserController) Index(c *fiber.Ctx) error {
	users, err := services.GetAllUsers(uc.db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	
	userResources := resources.TransformCollection(users)

	return c.JSON(userResources)
}
