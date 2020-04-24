package models

import "time"

type Dataset struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Owner       string     `json:"owner"`
	Size        int        `json:"size"`
	Features    int        `json:"features"`
	Bounds      [4]float64 `json:"bounds"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	Description string     `json:"description"`

	// gray add
	Source  string   `json:"-"`
	GeoType GeoType  `json:"geotype"`
	Fields  []string `json:"fields"`
}

// 保存dataset到数据库中
func (dt *Dataset)Save() error {
	return db.Create(dt).Error
}