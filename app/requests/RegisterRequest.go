package requests

import (
	"reflect"
	"tradeapi/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=3"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
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
			field := err.Field()
			jsonTag, ok := r.GetJSONTag(field)
			if ok {
				field = jsonTag
			}
			validationErrors[field] = getErrorMessage(err)
		}
	}

	if !r.ValidateEmailUniqueness() {
		validationErrors["email"] = "O e-mail já está em uso"
	}

	if len(validationErrors) > 0 {
		return map[string]interface{}{
			"error":     "validação",
			"details": validationErrors,
		}
	}

	return nil
}

func (r *RegisterRequestInternal) GetJSONTag(field string) (string, bool) {
	rv := reflect.ValueOf(r.RegisterRequest)
	for i := 0; i < rv.NumField(); i++ {
		if rv.Type().Field(i).Name == field {
			tag := rv.Type().Field(i).Tag.Get("json")
			if tag != "" {
				return tag, true
			}
		}
	}
	return "", false
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
	case "eqfield":
		return "Deve ser igual a " + err.Param()
	default:
		return "Entrada inválida"
	}
}
