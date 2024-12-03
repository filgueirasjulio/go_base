package resources

import (
	"tradeapi/app/models"
	"tradeapi/utils/helpers"
)

type UserResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone	  string  `json:"phone"`
	BirthDate string  `json:"birth_date"`
	CreatedAt string `json:"created_at"`
}

// Transform transforma um modelo User em um UserResource
func Transform(user *models.User) UserResource {
	return UserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone: helpers.FormatPhoneNumber(user.Phone),
		BirthDate: user.BirthDate.Format("02/01/2006"),
		CreatedAt: user.CreatedAt.Format("02/01/2006"),
	}
}

// TransformCollection transforma uma lista de modelos User em uma lista de UserResource
func TransformCollection(users []*models.User) []UserResource {
	resources := make([]UserResource, len(users))
	for i, user := range users {
		resources[i] = Transform(user)
	}
	return resources
}