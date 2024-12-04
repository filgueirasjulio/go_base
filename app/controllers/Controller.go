package controllers

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
	Logger zerolog.Logger
}

func NewController(db *gorm.DB, logger zerolog.Logger) *Controller {
	return &Controller{DB: db,
		Logger: logger,
	}
}