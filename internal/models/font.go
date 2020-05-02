package models

import "Mapscope/internal/database"

const (
	// PBFONTEXT pbf fonts package format
	PBFONTEXT = ".pbfonts"
	// DEFAULTFONT
	DEFAULTFONT = "Noto Sans Regular"
)

type Font struct {
	Id          string `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"index;not null"`
	Path        string `json:"path" gorm:"not null"`
	Owner       string `json:"owner" gorm:"index"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	URL         string `json:"url"`
	Public      int    `json:"public"`
	Compression int    `json:"compression"`
}

func (ft *Font) Save() error  {
	return database.Get().Save(ft).Error
}