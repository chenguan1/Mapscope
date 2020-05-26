package services

import (
	"Mapscope/internal/config"
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/thirdparty/teris-io/shortid"
	"Mapscope/internal/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 根据id获取dataset
func DatasetGet(dataset_id string) (*models.Dataset, error) {
	var dt models.Dataset
	err := database.Get().Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return nil, fmt.Errorf("DatasetGet failed, err: %v", err)
	}
	return &dt, nil
}

// 删除指定的dataset
func DatasetDelete(dataset_id string) error {
	var err error
	var dt models.Dataset

	db := database.Get()
	err = db.Where(models.Dataset{Id: dataset_id}).Find(&dt).Error
	if err != nil {
		return fmt.Errorf("DatasetDelete err: %v", err)
	}

	// 事务删除
	tx := db.Begin()

	// 1.删除dataset_did表
	if err = tx.DropTableIfExists(dt.TableName).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Database droptable &v err: %v", dt.TableName, err)
	}
	// 2.删除dataset记录
	if err = tx.Delete(models.Dataset{Id: dt.Id}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Database delete dataset record err: %v", err)
	}

	// 3.删除备份
	if err = DataBackupDeleteAll(dataset_id); err != nil {
		tx.Rollback()
		return fmt.Errorf("Database delete backups err: %v", err)
	}

	// 提交
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Database commit err: %v", err)
	}

	return nil
}

// 获取某用户的所有dataset
func DatasetList(user string) ([]models.Dataset, error) {
	var dts []models.Dataset
	err := database.Get().Where(models.Dataset{Owner: user}).Find(&dts).Error
	if err != nil {
		return nil, fmt.Errorf("DatasetList err: %v", err)
	}
	return dts, nil
}

// 从上传的文件中导入Dataset
func DatasetsFromUpload(files []models.FileUped, user string) ([]models.Dataset, error) {
	dtsrcs, err := DatasourceFromUpload(files, user)
	if err != nil {
		return nil, fmt.Errorf("DatasetsFromUpload failed, err: %v", err)
	}

	// 入库 dtsrc，并转换成dataset
	dts := make([]models.Dataset, 0)
	for _, ds := range dtsrcs {
		// 入库
		err := ds.Save()
		if err != nil {
			return nil, fmt.Errorf("DatasetsFromUpload, Datasorce save failed, path %v, err: %v", ds.Path, err)
		}
		// to dataset
		dt, err := Datasource2Dataset(&ds)
		if err != nil {
			return nil, fmt.Errorf("DatasetsFromUpload, Dataset save failed, path %v, err: %v", ds.Path, err)
		}

		// dt 入库
		err = dt.Save()
		if err != nil {
			return nil, fmt.Errorf("DatasetsFromUpload, Dataset insert failed, path %v, err: %v", ds.Path, err)
		}
		dts = append(dts, *dt)
	}

	return dts, nil
}

// dataset to mvt
func DatasetToMvtBuf(dataset_id string, zoom, x, y int) ([]byte, error) {
	dt, err := DatasetGet(dataset_id)
	if err != nil {
		return nil, fmt.Errorf("DatasetToMvtBuf err: %v", err)
	}

	minx, miny := utils.TileUl_4326(zoom, x, y)
	maxx, maxy := utils.TileUl_4326(zoom, x+1, y+1)

	fields := dt.Fields.Keys()
	fields = append([]string{"gid"}, fields...)
	flds := `"` + strings.Join(fields, `","`) + `"`

	sqltmp := `
WITH mvtgeom AS
	(SELECT ST_AsMVTGeom(ST_Transform(%v,3857),ST_TileEnvelope(%v,%v,%v)) AS geom, %v
	FROM %v 
	WHERE ST_Intersects(%v,ST_MakeEnvelope(%v,%v,%v,%v, 4326)))
	SELECT ST_AsMvt(mvtgeom.*,'%v') AS mvt
	FROM mvtgeom
`

	sql := fmt.Sprintf(sqltmp, "geom", zoom, x, y, flds, dt.TableName, "geom", minx, miny, maxx, maxy, dt.Name)
	//fmt.Println(sql)

	type geoItem struct {
		Mvt []byte `gorm:"column:mvt"`
	}

	var gi geoItem
	err = database.Get().Raw(sql).Scan(&gi).Error
	if err != nil {
		return nil, fmt.Errorf("DatasetToMvtBuf query tile failed, err: %v", err)
	}
	return gi.Mvt, nil
}

// tilejson of dataset
func DatasetTilejson(dataset_id string) (*models.Tilejson, error) {
	dt, err := DatasetGet(dataset_id)
	if err != nil {
		return nil, fmt.Errorf("DatasetTilejson err: %v", err)
	}

	return dt.ToTileJson(), nil
}

// 多个 dataset 的tilejson
func DatasetsTilejson(datasets_id []string) (*models.Tilejson, error) {

	var owner string
	var tj *models.Tilejson
	for _, dtid := range datasets_id {
		dt, err := DatasetGet(dtid)
		if err != nil {
			return nil, fmt.Errorf("datasets tile json err: dataset &v can not get.", dtid)
		}
		owner = dt.Owner
		if tj == nil {
			tj = dt.ToTileJson()
		} else {
			t := dt.ToTileJson()
			tj.Name = tj.Name + "," + t.Name
			tj.VectorLayers = append(tj.VectorLayers, t.VectorLayers...)
			tj.Bounds[0], tj.Bounds[1] = math.Min(tj.Bounds[0], t.Bounds[0]), math.Min(tj.Bounds[1], t.Bounds[1])
			tj.Bounds[2], tj.Bounds[3] = math.Max(tj.Bounds[2], t.Bounds[2]), math.Max(tj.Bounds[3], t.Bounds[3])
		}
	}

	url := fmt.Sprintf("http://localhost:8080/datasets/v1/%s/%s/{z}/{x}/{y}.mvt", owner, strings.Join(datasets_id, ","))
	tj.Tiles = []string{url}
	tj.UpdateCenter()

	return tj, nil
}

// 提交编辑，version+1，
func DatasetCommit(dtid string) (*models.Dataset, error) {
	dt, err := DatasetGet(dtid)
	if err != nil {
		return nil, fmt.Errorf("DatasetCommit err: %v", err)
	}
	if dt.Edited == 0 {
		return nil, fmt.Errorf("DatasetCommit err: no edited.")
	}

	dt.Edited = 0
	dt.Version++
	dt.Save()
	return dt, nil
}

// 恢复到以前备份的版本，根据版本号恢复
func DatasetRevertTo(dtid string, version int) (*models.Dataset, error) {
	//
	return nil, nil
}


// dataset => tileset todo ...
func Dataset2Tileset(dtid string) (*models.Task, error) {
	dt, err := DatasetGet(dtid)
	if err != nil {
		return nil, fmt.Errorf("Dataset2Tileset err: %v", err)
	}

	// tileset是否已经存在，如果已经存在，则返回失败
	{
		ts,err := TilesetGet(dtid)
		if err == nil && dt.Version == ts.Version{
			return nil, fmt.Errorf("Dataset2Tileset err: Tileset exist.")
		}
	}


	sid ,_:= shortid.GenerateLower()
	task := &models.Task{
		ID:        sid,
		Base:      dt.Id,
		Name:      dt.Name,
		Type:      models.TSPUBLIC, // "tspublic"
		Owner:     dt.Owner,
		Progress:  0,
		Status:    "processing",
		Error:     "",
		Pipe:      make(chan struct{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 异步导入
	TaskSet.Store(task.ID, task)
	task.Save()
	go func(dt *models.Dataset, task *models.Task) {
		defer func(task *models.Task) {
			task.Pipe <- struct{}{}
		}(task)

		gjson := "" // geojson 文件路径
		// 如果geojson文件存在，则直接使用
		if dt.Version == 0{
			ds, err := DatasourceGet(dt.Source)
			if err == nil{
				if _,err = os.Stat(ds.Path); !os.IsNotExist(err){
					gjson = ds.Path
				}
			}
		}
		// 如果geojson文件不存在，则从数据库中导出一个geojson文件
		gexist := false
		if gjson != "" {
			if utils.PathExist(gjson){
				gexist = true
			}
		}

		if !gexist{
			// 不存在则导出
			folder := config.PathDatasources(dt.Owner)
			gname := fmt.Sprintf("%v_%v.geojson", dt.Id, dt.Version)
			gjson = filepath.Join(folder,gname)

			para := utils.Pg2geojsonParams{
				Pghost:    viper.GetString("db.host"),
				Pgport:    viper.GetString("db.port"),
				Pguser:    viper.GetString("db.user"),
				Pgpswd:    viper.GetString("db.password"),
				Dbname:    viper.GetString("db.database"),
				Geojson:   gjson,
				TableName: dt.TableName,
			}

			if _,err = os.Stat(gjson); os.IsNotExist(err){
				err := utils.Pg2geojson(para)
				if err != nil {
					task.Status = "failed"
					task.Error = fmt.Errorf("Pg2geojson failed: %v", err).Error()
					return
				}
			}
		}
		task.Progress = 20
		task.Update()

		// geojson 有了，转为mbtiles
		mbfolder := config.PathTilesets(dt.Owner)
		utils.EnsurePathExist(mbfolder)
		mbfile := fmt.Sprintf("%v_%v.mbtiles", dt.Name, dt.Version)
		mb := filepath.Join(mbfolder,mbfile)

		err = utils.CreateMbtiles([]string{gjson}, dt.Name, mb, task)
		if err != nil{
			task.Status = "failed"
			task.Error = fmt.Errorf("Create mbtiles failed: %v",err).Error()
			return
		}

		ts,err := TilesetLoad(mb)
		if err != nil{
			task.Status = "failed"
			task.Error = fmt.Errorf("Loadd mbtiles failed: %v",err).Error()
			return
		}

		// 信息
		ts.Id = dt.Id
		ts.Dataset = dt.Id
		ts.Version = dt.Version
		ts.Owner = dt.Owner
		ts.Name = ts.Name

		err = ts.Save()
		if err != nil{
			task.Status = "failed"
			task.Error = fmt.Errorf("save tileset info failed: %v",err).Error()
			return
		}

		task.Progress = 100
		task.Status = "finished"
		task.Error = ""

	}(dt, task)

	go func(task *models.Task) {
		<-task.Pipe
		task.Update()
		TaskSet.Delete(task.ID)
		if task.Error != ""{
			log.Errorf("public tileset failed: %v", task.Error)
		}
	}(task)

	return task, nil
}