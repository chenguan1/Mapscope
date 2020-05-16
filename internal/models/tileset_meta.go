package models

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
	Tilejson     string        `json:"tilejson"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Format       string        `json:"format"`
	Description  string        `json:"description"`
	Version      string        `json:"version"`
	Tiles        []string      `json:"tiles"`
	Minzoom      int           `json:"minzoom"`
	Maxzoom      int           `json:"maxzoom"`
	Bounds       [4]float64    `json:"bounds"`
	Center       [3]float64    `json:"center"`
	VectorLayers []VectorLayer `json:"vector_layers"`
}

// 从metadata表中解析信息
func (tm *TilesetMeta)Parse(meta map[string]interface{}) error {
	panic("todo ...")
	return nil
}