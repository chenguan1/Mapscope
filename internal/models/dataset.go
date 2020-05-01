package models

import (
	"Mapscope/internal/database"
	"fmt"
	"time"
)

type Dataset struct {
	Id           string    `json:"id" gorm:"primary_key"`
	Name         string    `json:"name"`
	Owner        string    `json:"owner"`
	Size         int64     `json:"size"`
	FeatureCount int       `json:"feature_count"`
	Extent       Bound     `json:"extent" gorm:"type:json"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
	Description  string    `json:"description"`

	// gray add
	Source    string `json:"-"`
	TableName string `json:"-"`

	Public  int     `json:"public"` // 是否可以公开访问
	GeoType GeoType `json:"geotype"`
	Fields  Fields  `json:"fields" gorm:"type:json"` // , , ,
	Version int     `json:"version"`                  // 版本号 0 开始，初始版本号为0
	Edited  int     `json:"edited"`                   // 编辑的次数，编辑状态下不缓存，不走缓存
}

// 保存dataset到数据库中
func (dt *Dataset) Save() error {
	db := database.Get()
	return db.Create(dt).Error
}

// 得到tilejson
func (dt *Dataset) ToTileJson() *Tilejson {
	tj := NewTileJson()
	tj.Name = dt.Name
	tj.Description = dt.Description
	//tj.Version = strconv.Itoa(dt.Version)
	tj.Version = "0.0.0"
	tj.Bounds[0] = dt.Extent[0]
	tj.Bounds[1] = dt.Extent[1]
	tj.Bounds[2] = dt.Extent[2]
	tj.Bounds[3] = dt.Extent[3]
	url := fmt.Sprintf("http://localhost:8080/datasets/v1/%s/%s/{z}/{x}/{y}.mvt", dt.Owner, dt.Id)
	tj.Tiles = append(tj.Tiles, url)
	tj.VectorLayers = append(tj.VectorLayers, vectorLayer{
		Id: dt.Name,
	})
	return tj
}
