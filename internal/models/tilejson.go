package models

// https://github.com/mapbox/tilejson-spec/tree/master/2.2.0
type Tilejson struct {
	Tilejson     string        `json:"tilejson"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Version      string        `json:"version"`
	Tiles        []string      `json:"tiles"`
	Minzoom      int           `json:"minzoom"`
	Maxzoom      int           `json:"maxzoom"`
	Bounds       [4]float64    `json:"bounds"`
	Center       [3]float64    `json:"center"`
	VectorLayers []vectorLayer `json:"vector_layers"`
}

type vectorLayer struct {
	Id string `json:"id"`
}

func NewTileJson() *Tilejson {
	return &Tilejson{
		Tilejson: "2.2.0",
		Version:  "1.0.0",
		Minzoom:  4,
		Maxzoom:  30,
		Bounds:   [4]float64{-180, -85.05112877980659, 180, 85.0511287798066},
		Center:   [3]float64{0, 0, 4},
	}
}
