package services

import (
	"tradeapi/repositories"
	"gorm.io/gorm"
	"tradeapi/models"
)

func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	return repositories.GetAllUsers(db)
}
