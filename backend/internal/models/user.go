package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	WalletAddress string `json:"wallet_address" gorm:"uniqueIndex;not null"`
	Username      string `json:"username" gorm:"uniqueIndex;not null"`
	Role          string `json:"role" gorm:"not null"`
	Bio           string `json:"bio"`
	AvatarURL     string `json:"avatar_url"`
}
