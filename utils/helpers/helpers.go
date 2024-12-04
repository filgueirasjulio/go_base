package helpers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
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
func ErrorResponse(c *fiber.Ctx, status int, data map[string]interface{}, logger zerolog.Logger) error {
	err := c.Status(status).JSON(data)
	logger.Error().Err(err).Int("status", status).Interface("data", data)
	return c.Status(status).JSON(data)
}

// mesmo objetivo da ErrorReponse, mas recebe uma string
func ErrorResponseString(c *fiber.Ctx, status int, mensagem string, logger zerolog.Logger) error {
	err := errors.New(mensagem)
	logger.Error().Err(err).Int("status", status)
	return c.Status(status).JSON(fiber.Map{"error": mensagem})
}

// gera código de 4 digitos
func GenerateValidationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(9999-1000)+1000)
}

// retorna telefone com mascara
func FormatPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) != 11 {
		return phoneNumber
	}

	formatted := fmt.Sprintf("(%s)%s-%s",
		phoneNumber[:2],  // Area Code
		phoneNumber[2:7], // Prefix
		phoneNumber[7:])  // Line Number

	return formatted
}

// retira mascara de telefone
func RemovePhoneMask(phone string) string {
	regexp := regexp.MustCompile(`[^\d]`)
	return regexp.ReplaceAllString(phone, "")
}

// Função auxiliar para retornar string ou nil se o telefone for vazio
func SafeGetPhone(phone string) interface{} {
    if phone == "" {
        return nil
    }
    return phone
}

// Função auxiliar para retornar string ou nil se a data de nascimento for nil
func SafeGetDate(date *time.Time) interface{} {
    if date == nil {
        return nil
    }
    return date.Format("02/01/2006")
}

// função para exibição de mensagens de sucesso
func SuccessResponseString(c *fiber.Ctx, status int, mensagem string, logger zerolog.Logger, start time.Time) error {
	logger.Info().Msg(mensagem)
	elapsedTime := time.Since(start)
	logger.Info().Msgf("Tempo de execução: %s", elapsedTime)
	return c.Status(status).JSON(fiber.Map{"success": mensagem})
}

// função para exibição de mensagens de sucesso e retorno de datos
func SuccessResponseData(c *fiber.Ctx, status int, data map[string]interface{}, mensagem string, logger zerolog.Logger, start time.Time) error {
	logger.Info().Msg(mensagem)
	elapsedTime := time.Since(start)
	logger.Info().Msgf("Tempo de execução: %s", elapsedTime)

	dados := fiber.Map{"success": mensagem}
	for k, v := range data {
		dados[k] = v
	}

	return c.Status(status).JSON(dados)
}
