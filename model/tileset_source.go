package model

// A tileset source is raw geographic data
// formatted as line-delimited GeoJSON
type TilesetSource struct {
	Id    string `json:"id"`
	Files int    `json:"files"`

	FileSize   int64 `json:"file_size"`
	SourceSize int64 `json:"source_size"`

	Size     int64  `json:"size"` // same as SourceSize
	SizeNice string `json:"size_nice"`
}
