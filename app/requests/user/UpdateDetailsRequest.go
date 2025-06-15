package User

import (
	"base/app/requests"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserUpdateDetailsRequestParams struct {
    Name string `json:"name" validate:"required" example:"João Souza"`
    Phone string `json:"phone" validate:"required" example:"71991112222"`
	BirthDate  string `json:"birth_date" validate:"required" example:"10/01/2000"`
}

type UserRequest struct {
	UserUpdateDetailsRequestParams
	DB *gorm.DB
}

func NewUserUpdateDetailsRequest(db *gorm.DB) *UserRequest {
    return &UserRequest{DB: db}
}

func (r *UserRequest) Validate() map[string]interface{} {
    validationErrors := make(map[string]string)
	req := requests.NewRequestHelper()
    
    // Validação da Step1
    err := validator.New().Struct(&r.UserUpdateDetailsRequestParams)
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