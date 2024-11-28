package config

import (
	"errors"
	"os"
	"strings"
	"tradeapi/app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Chave secreta para assinar/verificar tokens
var jwtKey = []byte(os.Getenv("JWT_KEY"))

func Middleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Primeiro, verifica se o token JWT está presente
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token não fornecido",
			})
		}

		// Verifica se o header contém "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Formato do token inválido",
			})
		}

		tokenString := parts[1]

		// Valida e analisa o token JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Verifica o método de assinatura (HS256)
			if t.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("método de assinatura inválido")
			}

			// Chave secreta usada para validar o token
			secretKey := os.Getenv("JWT_KEY")
			if secretKey == "" {
				return nil, errors.New("chave secreta não configurada")
			}
			return []byte(secretKey), nil
		})

		// Verifica se o token é válido
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido: o token está malformado ou contém um número inválido de segmentos",
			})
		}

		// Armazena as claims no contexto para uso posterior
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Falha ao processar claims do token",
			})
		}

		// Obtém o user_id do token diretamente das claims
		userID := claims["user_id"].(string)

		// Verifica se o código de validação foi validado para este usuário antes de prosseguir
		if !isValidationCodeValid(userID, db) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Código de validação não foi validado",
			})
		}

		// Armazena o user_id no contexto para ser usado nas rotas seguintes
		c.Locals("user", userID)

		// Prossegue para o próximo handler
		return c.Next()
	}
}

// Função que verifica se o código de validação foi validado
func isValidationCodeValid(userID string, db *gorm.DB) bool {
	// Busca o último código de validação para o usuário
	var validationCode models.ValidationCode
	result := db.Where("user_id = ?", userID).Order("created_at desc").First(&validationCode)
	if result.Error != nil {
		return false
	}

	return validationCode.IsValidated
}
