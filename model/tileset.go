package model

type TilesetCreateForm struct {
	Rcp         Recipe `json:"recipe"`
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
}
