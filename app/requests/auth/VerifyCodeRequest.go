package auth

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
    "base/app/requests"
)

type VerifyCodeRequestParams struct {
    Email string `json:"email" validate:"required,email" example:"joao@example.com"`
}

type VerifyCodeRequest struct {
	VerifyCodeRequestParams
	DB *gorm.DB
}

func NewVerifyCodeRequest(db *gorm.DB) *VerifyCodeRequest {
    return &VerifyCodeRequest{DB: db}
}

func (l *VerifyCodeRequest) Validate() map[string]interface{} {
    validationErrors := make(map[string]string)
	req := requests.NewRequestHelper()
    
    // Validação da Step1
    err := validator.New().Struct(&l.VerifyCodeRequestParams)
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