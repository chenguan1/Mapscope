package models

import "time"

// 数据的编辑记录
// 根据此表可恢复
type Edition struct {
	Id         string    `json:"id" gorm:"primary_key"`
	Username   string    `json:"username"`
	DatasetId  string    `json:"dataset_id"`
	EditionNo  int       `json:"edition_no"`
	TableName  string    `json:"-"`
	EditedTime time.Time `json:"edited_time"`
}
