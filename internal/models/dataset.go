package models

import (
	"Mapscope/internal/database"
	"Mapscope/internal/thirdparty/paulmach/orb"
	"time"
)

type Dataset struct {
	Id          string     `json:"id" gorm:"primary_key"`
	Name        string     `json:"name"`
	Owner       string     `json:"owner"`
	Size        int64      `json:"size"`
	FeatureCount    int    `json:"feature_count"`
	Extent      orb.Bound `json:"extent" gorm:"type:json"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	Description string     `json:"description"`

	// gray add
	Source  string   `json:"-"`
	TableName string `json:"-"`

	GeoType GeoType  `json:"geotype"`
	Fields  string `json:"fields"` // , , ,
	EditionNo int `json:"edition_no"` // 版本号 0 开始
	Edited bool `json:"edited"` // 被编辑过，则需要提交版本更新，提交更新后，备份当前数据，并Version++
}

// 保存dataset到数据库中
func (dt *Dataset)Save() error {
	db := database.Get()
	return db.Create(dt).Error
}