package models

import "time"

// 数据的编辑记录
// 根据此表可恢复
/*
数据版本管理策略：
	edited == 0 为正常模式，edited > 0 为编辑状态，该状态下缓存失效
	编辑，如果edited == 0则先备份备份dataset数据表，edited = 1(不管编辑是否成功)
	提交编辑，如果edited != 0 {edited = 0; version++}
*/
type DataBackup struct {
	Id          int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Dataset     string    `json:"dataset"` // dataset id
	Source      string    `json:"source"` // datasource id // 可为空
	Version     int       `json:"version"`
	TableName   string    `json:"-"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
}
