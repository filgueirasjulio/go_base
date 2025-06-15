package resources

import (
	"github.com/iancoleman/orderedmap"
	"base/app/models"
	"base/utils/helpers"
)

// Transform transforma um modelo User em um UserResource
func Transform(user *models.User) *orderedmap.OrderedMap {
	om := orderedmap.New()
	if user != nil {
		om.Set("id", user.ID)
		om.Set("name", user.Name)
		om.Set("email", user.Email)
		om.Set("phone", helpers.SafeGetPhone(user.Phone))
		om.Set("birth_date", helpers.SafeGetDate(user.BirthDate))
		om.Set("created_at", user.CreatedAt.Format("02/01/2006"))
	}
	return om
}

// TransformCollection transforma uma lista de modelos User em uma lista de UserResource
func TransformCollection(users []*models.User) []*orderedmap.OrderedMap {
	resources := make([]*orderedmap.OrderedMap, len(users))
	for i, user := range users {
		resources[i] = Transform(user)
	}
	return resources
}
