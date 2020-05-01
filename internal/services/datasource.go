package services

import (
	"Mapscope/internal/config"
	"Mapscope/internal/models"
	"Mapscope/internal/utils"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 从上传的文件中解析出Datasource列表
func DatasourceFromUpload(files []models.FileUped, user string) ([]models.Datasource, error) {

	// 支持的所有后缀
	exts := models.GEOJSONEXT + models.SHPEXT + models.CSVEXT // + models.KMLEXT + models.GPXEXT

	// 得到所有文件
	allfs := make([]models.FileUped, 0)
	for _, ff := range files {
		switch ff.Ext {
		case models.ZIPEXT:
			fs, err := utils.UnzipFile(ff.Path, "")
			if err != nil {
				return nil, fmt.Errorf("DatasetUpload, unzip file failed, file %v, err: %v", ff.Path, err)
			}
			for _, vf := range fs {
				info, _ := os.Stat(vf)
				ext := strings.ToLower(filepath.Ext(info.Name()))
				if strings.Contains(exts, ext) {
					allfs = append(allfs, models.FileUped{
						Sid:  ff.Sid,
						Path: vf,
						Name: strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
						Ext:  strings.ToLower(filepath.Ext(info.Name())),
						Size: info.Size(),
					})
				}
			}

		case models.GEOJSONEXT:
			fallthrough
		case models.CSVEXT:
			allfs = append(allfs, ff)
		}
	}

	dsfolder := config.PathDatasources(user)
	utils.EnsurePathExist(dsfolder)

	// 将支持的格式转成geojson格式的datasorce
	dtsrcs := make([]models.Datasource, 0) // datasorce 准备入库，记录必要信息
	for _, vf := range allfs {
		gjname := vf.Sid + "_" + vf.Name + models.GEOJSONEXT
		gjpath := filepath.Join(dsfolder, gjname)

		err := ToGeojson(vf.Path, gjpath)
		if err != nil {
			return nil, fmt.Errorf("DatasourceFromUpload, convert to geojson failed: %v", err)
		}

		info, err := os.Stat(gjpath)
		if err != nil {
			return nil, fmt.Errorf("DatasourceFromUpload error: %v", err)
		}

		// to datasource
		st := models.Datasource{
			Id:        vf.Sid,
			Name:      vf.Name,
			Owner:     user,
			Src:       vf.Path,
			Path:      gjpath,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Size:      info.Size(),
			Tag:       "",
		}

		// append
		dtsrcs = append(dtsrcs, st)
	}

	return dtsrcs, nil
}

// 将datasorce（geojson）数据导入到数据库中
// 每个geojson成为一个数据表，数据表信息存入Dataset对象返回
// dataset对象不入库，让调用者处理
func Datasource2Dataset(ds *models.Datasource) (*models.Dataset, error) {
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
	ps.Srs_t = "EPSG:4326"
	ps.GeoName = "geom"
	ps.TableName = tableName

	// 入库
	err := utils.Ogr2Db(ds.Path, ps)
	if err != nil {
		return nil, fmt.Errorf("Datasource to database failed. err: %v", err)
	}

	// dataset 信息
	para := utils.OgrinfoPgParams{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DbName:   viper.GetString("db.database"),
	}

	info, err := utils.OgrinfoPg(para, tableName)
	if err != nil {
		return nil, fmt.Errorf("Get Dataset info failed, err :%v", err)
	}

	dt := &models.Dataset{
		Id:           sid,
		Name:         ds.Name,
		TableName:    tableName,
		Owner:        ds.Owner,
		Size:         ds.Size,
		FeatureCount: info.FeatureCount,
		Extent:       info.Extent,
		Created:      time.Now(),
		Modified:     time.Now(),
		Description:  "",
		Source:       ds.Id,
		GeoType:      models.GeoType(info.GeoType),
		Fields:       info.Fields,
	}

	return dt, nil
}
