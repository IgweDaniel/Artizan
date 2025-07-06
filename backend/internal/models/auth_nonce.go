package models

import "gorm.io/gorm"

type AuthNonce struct {
	gorm.Model
	WalletAddress string `json:"wallet_address" gorm:"uniqueIndex;not null"`
	Message       string `json:"message" gorm:"not null"`
}
