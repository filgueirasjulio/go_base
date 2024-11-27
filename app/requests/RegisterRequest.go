package requests

import (
	"tradeapi/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequestInternal struct {
	RegisterRequest
	DB *gorm.DB
}

func NewRegisterRequest(db *gorm.DB) *RegisterRequestInternal {
	return &RegisterRequestInternal{DB: db}
}

func (r *RegisterRequestInternal) Validate() map[string]interface{} {
	validationErrors := make(map[string]string)
	err := validator.New().Struct(&r.RegisterRequest)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = getErrorMessage(err)
		}
	}

	if !r.ValidateEmailUniqueness() {
		validationErrors["email"] = "O e-mail já está em uso"
	}

	if len(validationErrors) > 0 {
		return map[string]interface{}{
			"erro":     "validação",
			"detalhes": validationErrors,
		}
	}

	return nil
}

func (r *RegisterRequestInternal) ValidateEmailUniqueness() bool {
	var count int64
	r.DB.Model(&models.User{}).Where("email = ?", r.Email).Count(&count)

	return count == 0
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Este campo é obrigatório"
	case "email":
		return "Deve ser um e-mail válido"
	case "min":
		return "Deve ter pelo menos " + err.Param() + " caracteres"
	default:
		return "Entrada inválida"
	}
}
