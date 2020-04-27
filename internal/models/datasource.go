package models

import (
	"Mapscope/internal/database"
	"time"
)

// geojson格式的数据源
// 上传的数据为各种格式的原始格式
// 统一处理成相同格式的
type Datasource struct {
	Id        string    `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
	Src       string    `json:"src"` // 原始文件 上传到服务器的原始文件，存在upoads文件夹下
	Tag       string    `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// create or update datasource
// save to database, datasorces table
func (ds *Datasource) Save() error {
	db := database.Get()
	return db.Create(ds).Error
}


