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
	VectorLayers []VectorLayer `json:"vector_layers"`
}

func NewTileJson() *Tilejson {
	return &Tilejson{
		Tilejson: "2.2.0",
		Version:  "1.0.0",
		Minzoom:  4,
		Maxzoom:  24,
		Bounds:   [4]float64{-180, -85.05112877980659, 180, 85.0511287798066},
		Center:   [3]float64{0, 0, 6},
	}
}

func (tj *Tilejson) UpdateCenter() {
	long := (tj.Bounds[0] + tj.Bounds[2]) / 2.0
	lat := (tj.Bounds[1] + tj.Bounds[3]) / 2.0
	tj.Center = [3]float64{long, lat, float64(tj.Minzoom)}
}
