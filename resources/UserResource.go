package resources

import "tradeapi/models"

type UserResource struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// Transform transforma um modelo User em um UserResource
func Transform(user models.User) UserResource {
	return UserResource{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("02/01/2006"),
	}
}

// TransformCollection transforma uma lista de modelos User em uma lista de UserResource
func TransformCollection(users []models.User) []UserResource {
	resources := make([]UserResource, len(users))
	for i, user := range users {
		resources[i] = Transform(user)
	}
	return resources
}
