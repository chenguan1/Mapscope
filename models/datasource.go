package models

import (
	"Mapscope/global"
	"Mapscope/thirdparty/paulmach/orb"
	"Mapscope/utils"
	"fmt"
	"github.com/spf13/viper"
	"strings"
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
	db := global.GetDb()
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
	ps.Pghost = viper.GetString("db.host")
	ps.Pgport = viper.GetString("db.port")
	ps.Pguser = viper.GetString("db.user")
	ps.Pgpswd = viper.GetString("db.password")
	ps.Dbname = viper.GetString("db.database")
	ps.Srs_t =  "EPSG:4326"
	ps.GeoName = "geom"
	ps.TableName = tableName

	// 入库
	err := utils.Ogr2Db(ds.Path, ps)
	if err != nil{
		return nil, fmt.Errorf("Datasource to database failed. err: %v", err)
	}

	// dataset 信息
	para := utils.OgrinfoPgParams{
		Host:viper.GetString("db.host"),
		Port:viper.GetString("db.port"),
		Username:viper.GetString("db.user"),
		Password:viper.GetString("db.password"),
		DbName:viper.GetString("db.database"),
	}

	info ,err := utils.OgrinfoPg(para,tableName)
	if err != nil{
		return nil, fmt.Errorf("Get Dataset info failed, err :%v", err)
	}

	extend := orb.Bound{
		Min:orb.Point{info.Extent[0],info.Extent[1]},
		Max:orb.Point{info.Extent[2],info.Extent[3]},
	}

	dt := &Dataset{
		Id:sid,
		Name:info.LayerName,
		TableName:tableName,
		Owner:ds.Owner,
		Size:ds.Size,
		FeatureCount:info.FeatureCount,
		Extent:extend,
		Created:time.Now(),
		Modified:time.Now(),
		Description:"",
		Source:ds.Id,
		GeoType:GeoType(info.GeoType),
		Fields:strings.Join(info.Fields,","),
	}

	return dt, nil
}
