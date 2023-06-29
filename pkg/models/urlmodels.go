package models

import (
	"github.com/tolumadamori/scissor/pkg/config"
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	OriginalURL string
	CustomShort string
	UserID      uint
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
	CustomShort string `json:"custom_short"`
}

func (url *URL) Save() (*URL, error) {
	err := config.Database.Create(&url).Error
	if err != nil {
		return &URL{}, err
	}
	return url, nil
}
