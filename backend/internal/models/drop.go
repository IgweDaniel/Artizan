package models

import "gorm.io/gorm"

/*
id , collection_id , drop_type , start_time , price , supply
*/
type Drop struct {
	gorm.Model
	CollectionID uint        `json:"collection_id" gorm:"not null"`
	Collection   *Collection `json:"collection" gorm:"foreignKey:CollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DropType     string      `json:"drop_type" gorm:"not null"`  // e.g., "public", "private"
	StartTime    int64       `json:"start_time" gorm:"not null"` // Unix timestamp for the start time
	Price        int64       `json:"price" gorm:"not null"`      // Price in the smallest currency unit (e.g., cents for USD)
	Supply       int64       `json:"supply" gorm:"not null"`     // Total supply of NFTs in this drop
}
