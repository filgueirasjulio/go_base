package models

import (
	"time"
)

type RegisterHash struct {
    ID        uint   `gorm:"primary_key"`
    UserEmail string   `gorm:"not null"`
    Hash      string `gorm:"not null;size:32"`
	CreatedAt time.Time `gorm:"not null"`
    ExpiresAt time.Time `gorm:"not null"`
}