package utils

import (
	"fmt"
	"testing"
)

func TestConvertShp2Geojson(t *testing.T) {
	return
	err := ConvertShp2Geojson("./test/traffic.shp", "./test/traffic.geojson")
	if err != nil {
		t.Fail()
	}
}

func TestConvertCsv2Geojson(t *testing.T) {
	return
	err := ConvertCsv2Geojson("./test/0-2.csv", "./test/0.geojson")
	if err != nil {
		t.Fail()
	}
}

func TestOgr2Db(t *testing.T) {
	return
	para := NewOgr2DbParams()
	para.Pgport = "5432"
	para.Pghost = "localhost"
	para.Pguser = "postgres"
	para.Pgpswd = "111111"
	para.Dbname = "mapscope"
	para.TableName = "dataset_test0"

	in_json := "./test/0.geojson"
	err := Ogr2Db(in_json, para)
	if err != nil {
		t.Errorf("Ogr2Db Failed.err: %v", err)
	}
}

func TestOgrinfoPg(t *testing.T) {
	para := OgrinfoPgParams{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "111111",
		DbName:   "mapscope",
	}

	info, err := OgrinfoPg(para, "dataset_test0")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(info)

}
