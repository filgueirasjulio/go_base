package routes

import (
	"base/app/controllers"
	"base/app/middlewares"
	_ "base/utils/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Route struct {
	App *fiber.App
	DB  *gorm.DB
	Logger zerolog.Logger
}

// NewRoute cria uma nova instância de rotas com a aplicação e o banco de dados
func NewRoute(db *gorm.DB, app *fiber.App, logger zerolog.Logger) *Route {
	return &Route{App: app,
		DB: db,
		Logger: logger,
	}
}

// @title BaseGo
// @version 1.0
// @description API de exemplo para documentação com Swagger
// @termsOfService http://swagger.io/terms/
// @contact.name Suporte
// @contact.url http://www.basego.com/support
// @contact.email suporte@basego.com
// @license.name MIT
// @license.url http://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /

// SetupAPIRoutes configura as rotas da API
func (r *Route) SetupAPIRoutes() {
	// Inicializa o controller base com o banco de dados
	controllerBase := controllers.NewController(r.DB, r.Logger)

	// Define o prefixo da API
	api := r.App.Group("/api")

	/* Documentação com Swagger */
	api.Get("/swagger/*", fiberSwagger.New())

	/* Rotas de autenticação */
	authController := controllers.NewAuthController(controllerBase)
	auth := api.Group("/auth")
	{
		auth.Post("/login", authController.Login)
		auth.Post("/register_step_1", authController.RegisterStep1)
		auth.Post("/register_step_2", authController.RegisterStep2)
		auth.Post("/send-validation-code", authController.SendValidationCode)
		auth.Post("/verify-validation-code", authController.VerifyValidationCode)
		auth.Use(middlewares.Middleware(r.DB))
		{
			auth.Get("/logout", authController.Logout)
			auth.Post("/refresh-token", authController.HandleRefreshToken)
			auth.Post("/change-password", authController.RegisterStep2)
		}
	}

	/* Rotas de usuários */
	userController := controllers.NewUserController(controllerBase)
	users := api.Group("/users")
	users.Use(middlewares.Middleware(r.DB))
	{
		users.Get("/", userController.Index)
		users.Get("/:id", userController.Show)
		users.Put("/:id/update-details", userController.UpdateDetails)
	}
}
