package models

import (
	"Mapscope/internal/database"
	"time"
)

/*
Not Authorized - No Token	401	No token was used in the query.
Not Authorized - Invalid Token	401	Check the access token you used in the query.
This endpoint requires a token with {scope} scope	403	The access token used in the query needs the specified scope.
No such user	404	Check the username you used in the query.
Cannot find tileset	404	Check the tileset ID you used in the query.
The requested url's querystring \"limit\" property contains in invalid value.	422	The limit specified in the query is larger than 500, or contains non-numeric characters.
Invalid start key	422	Check the start key used in the query.
*/

// https://docs.mapbox.com/api/maps/#response-list-tilesets
type Tileset struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Format      string    `json:"format"`
	Filesize    int64     `json:"filesize"`
	Center      Center    `json:"center" gorm:"type:json"`
	Description string    `json:"description"`
	Visibility  string    `json:"visibility"`
	Public      int       `json:"public"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`

	Path     string      `json:"-"`
	Owner    string      `json:"owner"`
	Metadata TilesetMeta `json:"metadata" gorm:"type:json"`

	// mapscope
	Dataset string `json:"dataset"` // dataset id
	
	Version int    `json:"version"`
}

type TilesetCreateForm struct {
	Rcp         Recipe `json:"recipe"`
	Id          string `json:"id"` // same as name
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
}

type TilesetMetadata struct {
	Bounds       [4]float64
	Center       [3]float64
	Created      time.Time
	Filesize     int
	Format       string
	Id           string
	MapboxLogo   bool `json:"mapbox_logo"`
	Maxzoom      int
	Minzoom      int
	Modified     time.Time
	Name         string
	Private      bool
	Scheme       string
	Tilejson     string
	Tiles        []string
	VectorLayers []TilesetVectorLayer
	Version      string
	Webpage      string
}

// https://docs.mapbox.com/api/maps/#example-response-retrieve-tilejson-metadata
type TilesetVectorLayer struct {
	Description string
	Fields      map[string]string
}

// 保存dataset到数据库中,要保证Name 唯一
func (ts *Tileset) Save() error {
	db := database.Get()
	return db.Save(ts).Error
}
