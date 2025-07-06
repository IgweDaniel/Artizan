package models

import "gorm.io/gorm"

/*
nfts
id , drop_id , token_id , metadata_uri , lazy_minted
*/

type NFT struct {
	gorm.Model
	DropID      uint   `json:"drop_id" gorm:"not null"`
	Drop        *Drop  `json:"drop" gorm:"foreignKey:DropID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TokenID     string `json:"token_id" gorm:"not null;uniqueIndex"` // Unique identifier for the NFT
	MetadataURI string `json:"metadata_uri" gorm:"not null"`         // URI pointing to the NFT's metadata
	IsMinted    bool   `json:"is_minted" gorm:"default:false"`       // Indicates if the NFT is lazy minted
}
