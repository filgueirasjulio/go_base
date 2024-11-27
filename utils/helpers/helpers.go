package helpers

import (
	"strings"
	"github.com/gofiber/fiber/v2"
)

// A tag pode conter "name,omitempty" ou apenas "name". A função remove o omitempty
func ExtractJSONTag(tag string) string {
	
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag 
}

//função para exebição de erros
func ErrorResponse(c *fiber.Ctx, status int, data map[string]interface{}) error {
    return c.Status(status).JSON(data)
}

//mesmo objetivo da ErrorReponse, mas recebe uma string
func ErrorResponseString(c *fiber.Ctx, status int, mensagem string) error {
    return c.Status(status).JSON(fiber.Map{"erro": mensagem})
}
