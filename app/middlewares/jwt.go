package middlewares

import (
	"time"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Chave secreta para assinar/verificar tokens
var jwtKey = []byte(os.Getenv("JWT_KEY"))

func Middleware(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token não fornecido"})
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Formato do token inválido"})
        }

        tokenString := parts[1]
        claims := jwt.MapClaims{}
        _, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_KEY")), nil
        })

        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
        }

        if claims["exp"].(float64) < float64(time.Now().Unix()) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expirado"})
        }

        email := claims["email"].(string)
        c.Locals("user", email)
        return c.Next()
    }
}
