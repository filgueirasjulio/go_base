package helpers

import (
	"strings"
	"fmt"
	"math/rand"
	"time"
	"github.com/gofiber/fiber/v2"
)

// A tag pode conter "name,omitempty" ou apenas "name". A função remove o omitempty
func ExtractJSONTag(tag string) string {
	
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag 
}

//função para exibição de erros
func ErrorResponse(c *fiber.Ctx, status int, data map[string]interface{}) error {
    return c.Status(status).JSON(data)
}

//mesmo objetivo da ErrorReponse, mas recebe uma string
func ErrorResponseString(c *fiber.Ctx, status int, mensagem string) error {
    return c.Status(status).JSON(fiber.Map{"error": mensagem})
}

//função para exibição de mensagens de sucesso
func SuccessResponseString(c *fiber.Ctx, status int, mensagem string) error {
    return c.Status(status).JSON(fiber.Map{"success": mensagem})
}

//gera código de 4 digitos
func GenerateValidationCode() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%04d", rand.Intn(9999-1000)+1000)
}