package models

import (
	"Mapscope/utils"
	"fmt"
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
	return db.Create(ds).Error
}

// 将datasorce（geojson）数据导入到数据库中
// 每个geojson成为一个数据表，数据表信息存入Dataset对象返回
// dataset对象不入库，让调用者处理
func (ds *Datasource) ToDataset() (*Dataset, error)  {
	// 使用ogr2ogr工具，将geojson导入到数据库中

	// dataset_jalkdsf-dafd-d
	sid := ds.Id
	tableName := "dataset_" + sid

	// 参数
	ps := utils.NewOgr2DbParams()
	ps.Pghost = "localhost"
	ps.Pgport = "5432"
	ps.Pguser = "postgres"
	ps.Pgpswd = "123456"
	ps.Dbname = "mapscope"
	ps.Srs_t =  "EPSG:4326"
	ps.GeoName = "geom"
	ps.TableName = tableName

	// 入库
	err := utils.Ogr2Db(ds.Path, ps)
	if err != nil{
		return nil, fmt.Errorf("Datasource to database failed. err: %v", err)
	}

	// 类型

	dt := &Dataset{
		Id:sid,
		Name:ds.Name,
		Owner:ds.Owner,
		Size:ds.Size,
		Features:0,
		Bounds:[4]float64{0,0,0,0},
		Created:time.Now(),
		Modified:time.Now(),
		Description:"",
		Source:ds.Id,
		TableName:tableName,
		GeoType:"....",
		Fields:[]string{},
	}



	return dt, nil
}
