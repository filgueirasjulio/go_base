package helpers

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// converte formato da data
func ConvertData(data string) (time.Time, error) {
	layout := "02/01/2006"
	dataParseada, err := time.Parse(layout, data)
	if err != nil {
		log.Println(dataParseada)
		return dataParseada, err
	}
	return dataParseada, nil
}

// A tag pode conter "name,omitempty" ou apenas "name". A função remove o omitempty
func ExtractJSONTag(tag string) string {

	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

// função para exibição de erros
func ErrorResponse(c *fiber.Ctx, status int, data map[string]interface{}) error {
	return c.Status(status).JSON(data)
}

// mesmo objetivo da ErrorReponse, mas recebe uma string
func ErrorResponseString(c *fiber.Ctx, status int, mensagem string) error {
	return c.Status(status).JSON(fiber.Map{"error": mensagem})
}

//retorna telefone com mascara
func FormatPhoneNumber(phoneNumber string) string {
    if len(phoneNumber) != 11 {
        return phoneNumber
    }

    formatted := fmt.Sprintf("(%s)%s-%s",
        phoneNumber[:2], // Area Code
        phoneNumber[2:7], // Prefix
        phoneNumber[7:]) // Line Number

    return formatted
}

//retira mascara de telefone
func RemovePhoneMask(phone string) string {
	regexp := regexp.MustCompile(`[^\d]`) 
	return regexp.ReplaceAllString(phone, "")
}

// função para exibição de mensagens de sucesso
func SuccessResponseString(c *fiber.Ctx, status int, mensagem string) error {
	return c.Status(status).JSON(fiber.Map{"success": mensagem})
}

// gera código de 4 digitos
func GenerateValidationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(9999-1000)+1000)
}
