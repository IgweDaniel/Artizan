package models

import "gorm.io/gorm"

type Collection struct {
	gorm.Model
	CreatorID       uint   `json:"creator_id" gorm:"not null"`
	Creator         *User  `json:"creator" gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name            string `json:"name" gorm:"not null"`
	Description     string `json:"description" gorm:"not null"`
	ContractAddress string `json:"contract_address" gorm:"not null;uniqueIndex;type:varchar(42)"` // Ethereum address format
}
