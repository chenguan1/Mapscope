package model

type TilesetCreateForm struct {
	Rcp         Recipe `json:"recipe"`
	Id          string `json:"id"` // same as name
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
}
