package utils

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/teris-io/shortid"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ogr2dbparams struct {
	Pghost    string
	Pgport    string
	Pguser    string
	Pgpswd    string
	Dbname    string
	Srs_t     string
	GeoName   string
	TableName string
}

func NewOgr2DbParams() ogr2dbparams {
	sid, _ := shortid.Generate()
	return ogr2dbparams{
		Pghost:    "localhost",
		Pgport:    "5432",
		Pguser:    "postgres",
		Pgpswd:    "123456",
		Dbname:    "mapscope",
		Srs_t:     "EPSG:4326",
		GeoName:   "geom",
		TableName: "dataset_" + sid,
	}
}

func (para ogr2dbparams) ToStringList() []string {
	var pms []string
	pms = append(pms, []string{"-f", "PostgreSQL"}...)
	pg := fmt.Sprintf(`PG:dbname=%s host=%s port=%s user=%s password=%s`,
		para.Dbname, para.Pghost, para.Pgport, para.Pguser, para.Pgpswd)
	pms = append(pms, pg)
	pms = append(pms, []string{"-t_srs", para.Srs_t}...)
	//显示进度,读取outbuffer缓冲区
	pms = append(pms, "-progress")
	//跳过失败
	// pms = append(pms, "-skipfailures")//此选项不能开启，开启后windows会非常慢
	pms = append(pms, []string{"-nln", para.TableName}...)
	pms = append(pms, []string{"-lco", "FID=gid"}...)
	pms = append(pms, []string{"-lco", "GEOMETRY_NAME=" + para.GeoName}...)
	pms = append(pms, []string{"-lco", "LAUNDER=NO"}...)
	pms = append(pms, []string{"-lco", "EXTRACT_SCHEMA_FROM_LAYER_NAME=NO"}...)
	pms = append(pms, []string{"-lco", "OVERWRITE=YES"}...)
	pms = append(pms, []string{"-nlt", "PROMOTE_TO_MULTI"}...)
	return pms
}

// 调用ogr2ogr
func Ogr2Db(file_in string, p ogr2dbparams) error {
	file_in, err := filepath.Abs(file_in)
	if err != nil {
		return fmt.Errorf("file path is error: %v", err)
	}
	pms := p.ToStringList()
	pms = append(pms, file_in)

	if runtime.GOOS == "windows" {
		paramsString := strings.Join(pms, ",")
		decoder := mahonia.NewDecoder("gbk")
		paramsString = decoder.ConvertString(paramsString)
		pms = strings.Split(paramsString, ",")
	}

	cmd := exec.Command("ogr2ogr", pms...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("start insert to database failed. err: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait insert to database failed. err: %v", err)
	}

	return nil
}
