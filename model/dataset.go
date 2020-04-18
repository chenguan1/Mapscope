package model

import "time"

type Dataset struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Owner       string     `json:"owner"`
	Size        int        `json:"size"`
	Features    int        `json:"features"`
	Bounds      [4]float64 `json:"bounds"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	Description string     `json:"description"`
}
