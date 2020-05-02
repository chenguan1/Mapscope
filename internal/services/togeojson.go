package services

import (
	"Mapscope/internal/models"
	"Mapscope/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 将矢量数据转换成geojson格式文件
func ToGeojson(vfile string, out_geojson string) error {
	var err error
	if _, err = os.Stat(vfile); os.IsNotExist(err) {
		return err
	}

	if out_geojson == "" {
		out_geojson = strings.TrimSuffix(vfile, filepath.Ext(vfile)) + ".geojson"
	}

	ext := strings.ToLower(filepath.Ext(vfile))
	switch ext {
	case models.SHPEXT:
		err = utils.ConvertShp2Geojson(vfile, out_geojson)
		if err != nil {
			return err
		}
	case models.CSVEXT:
		err = utils.ConvertCsv2Geojson(vfile, out_geojson)
		if err != nil {
			return err
		}
	case models.GEOJSONEXT:
		_, err = utils.CopyFile(out_geojson, vfile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("file is not surpported, file: %v", vfile)
	}

	return nil
}
