package models

import "time"

type ValidationCode struct {
    ID           uint   `json:"id"`
    UserID       uint   `json:"user_id"`
    Code         string `json:"validation_code"`
	IsValidated bool `gorm:"default:false;not null" json:"is_validated"`
    CreatedAt    time.Time `json:"created_at"`
    ExpiresAt    time.Time `json:"expires_at"`
}