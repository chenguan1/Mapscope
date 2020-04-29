package services

import (
	"Mapscope/internal/database"
	"Mapscope/internal/models"
	"Mapscope/internal/utils"
	"fmt"
	"strings"
)

// 从上传的文件中导入Dataset
func DatasetsFromUpload(files []models.FileUped, user string) ([]models.Dataset, error)  {
	dtsrcs, err := DatasourceFromUpload(files, user)
	if err != nil{
		return nil, fmt.Errorf("DatasetsFromUpload failed, err: %v", err)
	}

	// 入库 dtsrc，并转换成dataset
	dts := make([]models.Dataset,0)
	for _,ds := range dtsrcs{
		// 入库
		err := ds.Save()
		if err != nil{
			return nil, fmt.Errorf("DatasetsFromUpload, Datasorce save failed, path %v, err: %v",ds.Path,err)
		}
		// to dataset
		dt,err := Datasource2Dataset(&ds)
		if err != nil{
			return nil, fmt.Errorf("DatasetsFromUpload, Dataset save failed, path %v, err: %v",ds.Path,err)
		}

		// dt 入库
		err = dt.Save()
		if err != nil{
			return nil, fmt.Errorf("DatasetsFromUpload, Dataset insert failed, path %v, err: %v",ds.Path,err)
		}
		dts = append(dts, *dt)
	}

	return dts, nil
}

// dataset to mvt
func DatasetToMvtBuf(dataset_id string, zoom, x, y int) ([]byte, error) {
	var dt models.Dataset
	err := database.Get().Where(models.Dataset{Id:dataset_id}).Find(&dt).Error
	if err != nil{
		return nil, fmt.Errorf("DatasetToMvtBuf failed, err: %v", err)
	}

	minx, miny := utils.TileUl_4326(zoom,x,y)
	maxx, maxy := utils.TileUl_4326(zoom,x+1,y+1)

	fields := strings.Split(dt.Fields,",")
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

	sql := fmt.Sprintf(sqltmp,"geom",zoom, x, y ,flds,dt.TableName,"geom",minx,miny,maxx,maxy,dt.Id)
	//fmt.Println(sql)

	type geoItem struct {
		Mvt []byte `gorm:"column:mvt"`
	}

	var gi geoItem
	err = database.Get().Raw(sql).Scan(&gi).Error
	if err != nil{
		return nil, fmt.Errorf("DatasetToMvtBuf query tile failed, err: %v", err)
	}
	return gi.Mvt, nil
}
