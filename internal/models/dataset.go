package models

import (
	"Mapscope/internal/database"
	"fmt"
	"time"
)

type Dataset struct {
	Id           string    `json:"id" gorm:"primary_key"`
	Name         string    `json:"name" gorm:"not null;unique"`
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
	Version int     `json:"version"`                 // 版本号 0 开始，初始版本号为0
	Edited  int     `json:"edited"`                  // 编辑的次数，编辑状态下不缓存，不走缓存
}

// 保存dataset到数据库中,要保证Name 唯一
func (dt *Dataset) Save() error {
	db := database.Get()

	// 检查是否重名，如果重名则需要将名称加上id
	var dtSameName []Dataset
	err := db.Where(Dataset{Name:dt.Name}).Find(&dtSameName).Error
	if err != nil{
		return fmt.Errorf("database cnn failed. %v", err)
	}

	// 如果重名，name加上sid
	if len(dtSameName) > 0{
		dt.Name = dt.Name + "_" + dt.Id
	}

	return db.Save(dt).Error
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

	long := (tj.Bounds[0] + tj.Bounds[2]) / 2.0
	lat  := (tj.Bounds[1] + tj.Bounds[3]) / 2.0
	tj.Center = [3]float64{long,lat,float64(tj.Minzoom)}

	return tj
}


// 判断是否处于编辑状态
func (dt *Dataset)IsEditing() bool  {
	return dt.Edited != 0
}