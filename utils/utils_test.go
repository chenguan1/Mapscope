package utils

import "testing"


func TestConvertShp2Geojson(t *testing.T) {
	err := ConvertShp2Geojson("./test/traffic.shp", "./test/traffic.geojson")
	if err != nil{
		t.Fail()
	}
}

func TestConvertCsv2Geojson(t *testing.T) {
	err := ConvertCsv2Geojson("./test/0-2.csv", "./test/0.geojson")
	if err != nil{
		t.Fail()
	}
}