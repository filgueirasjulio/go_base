// repositories/userRepository.go
package repositories

import (
	"tradeapi/models"
	"gorm.io/gorm"
)

var db *gorm.DB 

func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
