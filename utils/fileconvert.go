package utils

import (
	"Mapscope/pkg/csv2geojson"
	"fmt"
	"github.com/cartosquare/shp2geojson"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConvertShp2Geojson(shpfile string, geojsonfile string) error {
	data, err := shp2geojson.Convert(shpfile)
	if err != nil{
		return fmt.Errorf("convert shp to geojson failed: %v",err)
	}

	EnsurePathExist(filepath.Dir(geojsonfile))
	err = ioutil.WriteFile(geojsonfile,data,os.ModePerm)
	if err != nil{
		return fmt.Errorf("create geojson file failed, err: %v",err)
	}

	return nil
}

func ConvertCsv2Geojson(csvfile string, geojsonfile string) error {
	err := csv2geojson.Convert(csvfile, geojsonfile)
	if err != nil{
		return fmt.Errorf("convert csv to geojson failed: %v",err)
	}
	return  nil
}
