package models

// 上传的文件
type FileUped struct {
	Sid        string // 分配的shortid
	Path       string // 保存的路径
	Name       string // 去除文件名后的名字
	Ext        string // .zip .geojson
	Size       int64
}
