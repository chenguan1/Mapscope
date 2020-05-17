package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

/*
name	/data/mbtiles/berlin-big-summary2.mbtiles
description	/data/mbtiles/berlin-big-summary2.mbtiles
version	2
minzoom	0
maxzoom	14
center	11.909180,53.357101,14
bounds	11.909180,51.645286,14.721680,53.357101
type	overlay
format	pbf
json	{"vector_layers": [ { "id": "berlingbigsummarygeojson", "description": "", "minzoom": 0, "maxzoom": 14, "fields": {"oCount": "Number", "ratio": "Number", "uCount": "Number"} } ] }
*/

// metadata 表信息
type TilesetMeta struct {
	Tilejson    string `json:"tilejson"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Description string `json:"description"`
	//	Version      string        `json:"version"`
	Minzoom      int           `json:"minzoom"`
	Maxzoom      int           `json:"maxzoom"`
	Bounds       [4]float64    `json:"bounds"`
	Center       [3]float64    `json:"center"`
	VectorLayers []VectorLayer `json:"vector_layers"`
}

func (b TilesetMeta) Value() (value driver.Value, err error) {

	data, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

func (b *TilesetMeta) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Invalid Scan Source")
	}

	return json.Unmarshal(s, b)
}

// 从metadata表中解析信息
func (tm *TilesetMeta) Parse(meta map[string]interface{}) error {
	tm.Tilejson = "2.1.0"
	for k, v := range meta {
		switch k {
		case "name":
			tm.Name = v.(string)
		case "minzoom":
			tm.Minzoom = v.(int)
		case "maxzoom":
			tm.Maxzoom = v.(int)
		case "type":
			tm.Type = v.(string)
		case "format":
			tm.Format = v.(string)
		case "description":
			tm.Description = v.(string)
		case "center":
			items := v.([]float64)
			if len(items) != 3 {
				continue
			}
			tm.Center[0] = items[0]
			tm.Center[1] = items[1]
			tm.Center[2] = items[2]
		case "bounds":
			items := v.([]float64)
			if len(items) != 4 {
				continue
			}
			tm.Bounds[0] = items[0]
			tm.Bounds[1] = items[1]
			tm.Bounds[2] = items[2]
			tm.Bounds[3] = items[3]
		case "json":
			vlObj := v.(map[string]interface{})
			vls, ok := vlObj["vector_layers"]
			if !ok {
				log.Error("cannot read vector_layers info from metadata table.")
				continue
			}
			tm.VectorLayers = make([]VectorLayer, 0)
			err := json.Unmarshal([]byte(vls.(string)), &tm.VectorLayers)
			if err != nil {
				log.Errorf("parse vector_layer info failed: %v", err)
				return fmt.Errorf("parse vector_layer info failed: %v", err)
			}
		}
	}

	return nil
}
