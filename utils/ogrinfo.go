package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OgrinfoPgParams struct {
	Host string
	Port string
	DbName string
	Username string
	Password string
}

type OgrInfo struct {
	LayerName string
	GeoType string
	FeatureCount int
	Extent [4]float64
	Wkt string
	FidColumn string
	GeoColumn string
	Fields []string
}


func OgrinfoPg(p OgrinfoPgParams, layerName string) (*OgrInfo, error) {
	//ogrinfo -so "PG:dbname=mapscope host=localhost port=5432 user=postgres password=111111" dataset_test0

	var pms []string
	pg := fmt.Sprintf(`PG:dbname=%s host=%s port=%s user=%s password=%s`,
		p.DbName, p.Host, p.Port, p.Username, p.Password)
	pms = append(pms, "-so")
	pms = append(pms, pg)
	pms = append(pms, layerName)

	cmd := exec.Command("ogrinfo", pms...)

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
		return nil, fmt.Errorf("get dataset info failed. err: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("wait dataset info failed. err: %v", err)
	}
	time.Sleep(time.Millisecond * 100)
	rawinfo := stdoutBuf.String()

	ogrinfo,err := parseOgrinfo(rawinfo)
	if err != nil{
		return nil, fmt.Errorf("parse dataset info failed. err: %v", err)
	}

	return ogrinfo, nil
}

func parseOgrinfo(str string) (*OgrInfo, error) {
	layerName := regexp.MustCompile(`Layer name: .+`).FindString(str)
	layerName = strings.TrimPrefix(layerName,"Layer name: ")
	//fmt.Println(layerName)

	geotype := regexp.MustCompile(`Geometry: .+`).FindString(str)
	geotype = strings.TrimPrefix(geotype,"Geometry: ")
	geotype = strings.TrimSpace(geotype)
	//fmt.Println(geotype)

	features := regexp.MustCompile(`Feature Count: .+`).FindString(str)
	features = strings.TrimPrefix(features,"Feature Count: ")
	features = strings.TrimSpace(features)
	fmt.Println(features)
	featureCount,_ := strconv.Atoi(features)
	//fmt.Println(featureCount)

	extent := regexp.MustCompile(`Extent: .+`).FindString(str)
	extent = strings.TrimPrefix(extent,"Extent: ")
	extent = strings.TrimSpace(extent)

	nums := regexp.MustCompile(`-?[0-9]+\.*[0-9]*`).FindAllString(extent,-1)
	if len(nums) != 4{
		return nil, fmt.Errorf("cannot get extent info.")
	}
	var et [4]float64
	for i,v := range nums{
		et[i],_ = strconv.ParseFloat(v,64)
	}
	//fmt.Println(et)

	wkt := regexp.MustCompile(`Layer SRS WKT:[\s\S]+]\s`).FindString(str)
	wkt = strings.TrimSpace(strings.TrimPrefix(wkt,"Layer SRS WKT:"))
	//fmt.Println(wkt)

	fidColumn := regexp.MustCompile(`FID Column = .+`).FindString(str)
	fidColumn = strings.TrimPrefix(fidColumn,"FID Column = ")
	fidColumn = strings.TrimSpace(fidColumn)
	//fmt.Println(fidColumn)

	geoColumn := regexp.MustCompile(`Geometry Column = .+`).FindString(str)
	geoColumn = strings.TrimPrefix(geoColumn,"Geometry Column = ")
	geoColumn = strings.TrimSpace(geoColumn)
	//fmt.Println(geoColumn)

	fds := make([]string,0)
	filedsStrs := regexp.MustCompile(`.+: \w+ \([.0-9]+\)`).FindAllString(str,-1)
	for _,s := range filedsStrs{
		s = s[:strings.Index(s,":")]
		fds = append(fds,s)
	}

	//fmt.Println(fds)

	info := OgrInfo{
		LayerName: layerName,
		GeoType: geotype,
		FeatureCount: featureCount,
		Extent : et,
		Wkt : wkt,
		FidColumn : fidColumn,
		GeoColumn : geoColumn,
		Fields :fds,
	}

	//fmt.Printf("%#v",info)


	return &info, nil
}