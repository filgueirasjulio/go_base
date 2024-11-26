package services

import (
	"tradeapi/app/repositories"
	"gorm.io/gorm"
	"tradeapi/app/models"
)

func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	return repositories.GetAllUsers(db)
}
