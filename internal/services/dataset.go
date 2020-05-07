package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/utils"
	"fmt"
	"math"
	"strings"
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
func DatasetsTilejson(datasets_id []string) (*models.Tilejson, error)  {

	var owner string
	var tj *models.Tilejson
	for _, dtid := range datasets_id{
		dt,err := DatasetGet(dtid)
		if err != nil{
			return nil, fmt.Errorf("datasets tile json err: dataset &v can not get.", dtid)
		}
		owner = dt.Owner
		if tj == nil{
			tj = dt.ToTileJson()
		}else{
			t := dt.ToTileJson()
			tj.Name = tj.Name + ","+t.Name
			tj.VectorLayers = append(tj.VectorLayers, t.VectorLayers...)
			tj.Bounds[0], tj.Bounds[1] = math.Min(tj.Bounds[0], t.Bounds[0]), math.Min(tj.Bounds[1], t.Bounds[1])
			tj.Bounds[2], tj.Bounds[3] = math.Max(tj.Bounds[2], t.Bounds[2]), math.Max(tj.Bounds[3], t.Bounds[3])
		}
	}

	url := fmt.Sprintf("http://localhost:8080/datasets/v1/%s/%s/{z}/{x}/{y}.mvt", owner, strings.Join(datasets_id,","))
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
	if dt.Edited == 0{
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