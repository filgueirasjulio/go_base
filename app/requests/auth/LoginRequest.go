package auth

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
    "tradeapi/app/requests"
)

type LoginRequestParams struct {
    Email string `json:"email" validate:"required,email" example:"joao@example.com"`
	Password  string `json:"password" validate:"required,min=6" example:"senha123"`
}

type LoginRequest struct {
	LoginRequestParams
	DB *gorm.DB
}

func NewLoginRequest(db *gorm.DB) *LoginRequest {
    return &LoginRequest{DB: db}
}

func (l *LoginRequest) Validate() map[string]interface{} {
    validationErrors := make(map[string]string)
	req := requests.NewRequestHelper()
    
    // Validação da Step1
    err := validator.New().Struct(&l.LoginRequestParams)
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
    
    if len(validationErrors) > 0 {
        return map[string]interface{}{
            "error":     "validação",
            "details": validationErrors,
        }
    }
    
    return nil
}