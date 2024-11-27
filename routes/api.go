package routes

import (
	"tradeapi/app/controllers"
	_ "tradeapi/utils/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
	"gorm.io/gorm"
)

type Route struct {
	App *fiber.App
	DB  *gorm.DB
}

func NewRoute(db *gorm.DB, app *fiber.App) *Route {
	return &Route{App: app, DB: db}
}

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
func (r *Route) SetupAPIRoutes() {
	controllerBase := controllers.NewController(r.DB)

	api := r.App.Group("/api")

	api.Get("/swagger/*", fiberSwagger.New())

	/*usuários*/
	userController := controllers.NewUserController(controllerBase)
	users := api.Group("/users")

	users.Get("/", userController.Index)

	/*autenticação*/
	authController := controllers.NewAuthController(controllerBase)
	auth := api.Group("/auth/")

	auth.Post("/register", authController.Register)

	//listar rotas
	r.App.Get("/routes", func(c *fiber.Ctx) error {
		return c.JSON(r.App.Stack())
	})
}
