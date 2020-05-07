package utils

import (
	"math"
)

// tile 号转经纬度
func TileUl_4326(zoom, x, y int) (lon, lat float64) {
	n := math.Pi - 2.0*math.Pi*float64(y)/math.Exp2(float64(zoom))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	lon = float64(x)/math.Exp2(float64(zoom))*360.0 - 180.0
	return
}

// 经纬度转 tile号
func TileDeg2num(long, lat float64, zoom int) (y,x int) {
	x = int(math.Floor((long + 180.0) / 360.0 * (math.Exp2(float64(zoom)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(zoom)))))
	return
}