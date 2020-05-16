package services

import (
	"Mapscope/internal/models"
)

// 计算中心
func TileJsonCalCenter(tj *models.Tilejson) [3]float64 {
	long := (tj.Bounds[0] + tj.Bounds[2]) / 2.0
	lat := (tj.Bounds[1] + tj.Bounds[3]) / 2.0
	return [3]float64{long, lat, float64(tj.Minzoom)}
}
