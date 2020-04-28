package utils

import "math"

// tile 号转经纬度
func TileUl_4326(zoom, x, y int) (lon, lat float64) {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(zoom))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lon = float64(x)/math.Exp2(float64(zoom))*360.0 - 180.0
	return
}
