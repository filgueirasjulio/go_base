package auth

import (
	"base/app/models"
    "base/app/requests"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterRequestStep1 struct {
    Name  string `json:"name" validate:"required,min=3" example:"João Silva"`
    Email string `json:"email" validate:"required,email" example:"joao@example.com"`
}
type RegisterRequestStep2 struct {
    Hash            string `json:"hash" validate:"required" example:"72fd3951e5bd758b7619ddc62af1f052"`
    Password        string `json:"password" validate:"required,min=6" example:"senha123"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" example:"senha123"`
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
    req := requests.NewRequestHelper()

    // Validação da Step1
    err := validator.New().Struct(&r.RegisterRequestStep1)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            field := err.Field()
            jsonTag, ok := req.GetJSONTag(field)
            if ok {
                field = jsonTag
            }
            validationErrors[field] = requests.GetErrorMessage(err)
        }
    }
    
    // Validação da Step2
    err = validator.New().Struct(&r.RegisterRequestStep2)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            field := err.Field()
            jsonTag, ok := req.GetJSONTag(field)
            if ok {
                field = jsonTag
            }
            validationErrors[field] = requests.GetErrorMessage(err)
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

func (r *RegisterRequestInternal) ValidateEmailUniqueness() bool {
	var count int64
	r.DB.Model(&models.User{}).Where("email = ?", r.Email).Count(&count)

	return count == 0
}

