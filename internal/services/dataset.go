package services

import (
	"Mapscope/internal/models"
	"fmt"
)

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