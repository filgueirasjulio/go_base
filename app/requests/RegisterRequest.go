package requests

import (
	"reflect"
	"tradeapi/app/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterRequestStep1 struct {
    Name  string `json:"name" validate:"required,min=3"`
    Email string `json:"email" validate:"required,email"`
}

type RegisterRequestStep2 struct {
	Hash string `json:"hash" validate:"required"`
    Password string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
type RegisterRequestInternal struct {
	RegisterRequestStep1
	RegisterRequestStep2
	DB *gorm.DB
}

func NewRegisterRequest(db *gorm.DB) *RegisterRequestInternal {
    return &RegisterRequestInternal{DB: db}
}

func (r *RegisterRequestInternal) Validate() map[string]interface{} {
    validationErrors := make(map[string]string)
    
    // Validação da Step1
    err := validator.New().Struct(&r.RegisterRequestStep1)
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
    
    // Validação da Step2
    err = validator.New().Struct(&r.RegisterRequestStep2)
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
    // Reflete RegisterRequestStep1
    rv1 := reflect.ValueOf(r.RegisterRequestStep1)
    for i := 0; i < rv1.NumField(); i++ {
        if rv1.Type().Field(i).Name == field {
            tag := rv1.Type().Field(i).Tag.Get("json")
            if tag != "" {
                return tag, true
            }
        }
    }

    // Reflete RegisterRequestStep2
    rv2 := reflect.ValueOf(r.RegisterRequestStep2)
    for i := 0; i < rv2.NumField(); i++ {
        if rv2.Type().Field(i).Name == field {
            tag := rv2.Type().Field(i).Tag.Get("json")
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
