package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
    _ "tradeapi/utils/docs" 
	"tradeapi/app/controllers"
	fiberSwagger "github.com/gofiber/swagger"
)

// @title TradeAPI
// @version 1.0
// @description API de exemplo para documentação com Swagger
// @termsOfService http://swagger.io/terms/
// @contact.name Suporte
// @contact.url http://www.tradeapi.com/support
// @contact.email suporte@tradeapi.com
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /
// SetupAPIRoutes configura as rotas da API
func SetupAPIRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")

	api.Get("/swagger/*", fiberSwagger.New())

	/*usuários*/
	userController := controllers.NewUserController(db)
	users := api.Group("/users")

	users.Get("/", userController.Index)
}
