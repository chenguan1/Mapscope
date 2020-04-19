package model

//https://docs.mapbox.com/help/troubleshooting/tileset-recipe-reference

type Recipe struct {
	Version int                    `json:"version"`
	Layers  map[string]RecipeLayer `json:"layers"`
}

type RecipeLayer struct {
	Source  string `json:"source"`
	Minzoom int    `json:"minzoom"`
	Maxzoom int    `json:"maxzoom"`

	// optional
	// Features interface{} `json:"features"`
	// Tiles    interface{} `json:"tiles"`
}
