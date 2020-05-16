package models

/*
{
"vector_layers":[
{
"id":"berlingbigsummarygeojson",
"description":"",
"minzoom":0,
"maxzoom":14,
"fields":{
"oCount":"Number",
"ratio":"Number",
"uCount":"Number"
}
}
]
}
*/

type VectorLayer struct {
	Id          string            `json:"id"`
	Minzoom     int               `json:"minzoom"`
	Maxzoom     int               `json:"maxzoom"`
	Description string            `json:"description"`
	Fields      map[string]string `json:"fields"`
}
