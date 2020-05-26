package services

import "sync"

var (
	TaskSet sync.Map // 任务列表，不存数据库
)
