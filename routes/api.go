package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
    _ "tradeapi/docs" 
	"tradeapi/controllers"
	fiberSwagger "github.com/gofiber/swagger"
)

// SetupAPIRoutes configura as rotas da API
func SetupAPIRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")

	api.Get("/swagger/*", fiberSwagger.New())

	/*usuários*/
	userController := controllers.NewUserController(db)
	users := api.Group("/users")

	users.Get("/", userController.Index)
}
