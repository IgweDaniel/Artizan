package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type StoreFront struct {
	gorm.Model
	UserID      uint        `json:"user_id" gorm:"not null"`
	User        *User       `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Slug        string      `json:"slug" gorm:"uniqueIndex;not null"`
	BannerURL   string      `json:"banner_url" gorm:"not null"`
	ThemeConfig ThemeConfig `json:"theme_config" gorm:"type:jsonb"`
}

type ThemeConfig map[string]interface{}

// Scan implements the Scanner interface.
func (r *ThemeConfig) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	// Convert the value to a byte slice.
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Unmarshal the JSON into the map.
	return json.Unmarshal(bytes, r)
}

// Value implements the Valuer interface.
func (r ThemeConfig) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}

	// Marshal the map into JSON.
	return json.Marshal(r)
}
