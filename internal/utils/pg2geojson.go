package utils

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Pg2geojsonParams struct {
	Pghost    string
	Pgport    string
	Pguser    string
	Pgpswd    string
	Dbname    string
	Geojson   string
	TableName string
}

func NewPg2geojsonParams() Pg2geojsonParams {
	return Pg2geojsonParams{
		Pghost:    "localhost",
		Pgport:    "5432",
		Pguser:    "postgres",
		Pgpswd:    "123456",
		Dbname:    "mapscope",
		Geojson:   "out.geojson",
		TableName: "dataset_in" ,
	}
}

func (para Pg2geojsonParams) ToStringList() []string {
	var pms []string
	pms = append(pms, []string{"-f", "geojson"}...)
	pms = append(pms, para.Geojson)
	pg := fmt.Sprintf("PG:dbname=%s host=%s port=%s user=%s password=%s",
		para.Dbname, para.Pghost, para.Pgport, para.Pguser, para.Pgpswd)
	pms = append(pms, pg)
	pms = append(pms, "-overwrite") // 清空重写
	//显示进度,读取outbuffer缓冲区
	pms = append(pms, "-progress")
	pms = append(pms, para.TableName)
	return pms
}

// 调用 ogr2ogr，将postgis 转为 geojson
func Pg2geojson(p Pg2geojsonParams) error {
	pms := p.ToStringList()
	if runtime.GOOS == "windows" {
		paramsString := strings.Join(pms, ",")
		decoder := mahonia.NewDecoder("gbk")
		paramsString = decoder.ConvertString(paramsString)
		pms = strings.Split(paramsString, ",")
	}

	fmt.Println(pms)

	cmd := exec.Command("ogr2ogr", pms...)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	go func() {
		io.Copy(stdout, stdoutIn)
	}()
	go func() {
		io.Copy(stderr, stderrIn)
	}()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("start convert pg to geojson failed. err: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait convert pg to geojson failed. err: %v", err)
	}

	return nil
}